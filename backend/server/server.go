package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"net/textproto"
	"os"
	go_runtime "runtime"
	"strings"
	"time"

	"github.com/XHXHXHX/medical_marketing/errs"
	grpcUtil "github.com/XHXHXHX/medical_marketing/util/grpc"
	"github.com/XHXHXHX/medical_marketing/util/logx"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Server struct {
	addr string

	DailTimeout     time.Duration // 空则默认 5s
	ShutdownTimeout time.Duration // 空则默认 20s

	ServerRegisters []ServerRegister

	// 加载自定义的 http handler, 例如
	// func(mux *http.ServeMux) {
	//    mux.Handle("/abc/", myhandler)
	// }
	CustromHTTPHandlerSetFunc func(mux *http.ServeMux)

	// 可选, 增加自定义的 UnaryServerInterceptor
	UnaryServerInterceptors []grpc.UnaryServerInterceptor

	// 可选, 如果不指定,使用默认处理方式.
	HTTPErrorHandleFunc HTTPErrorHandleFunc

	// 允许透进来的 HTTP headers
	IncomingHeaderWhiteList []string

	DisableHeathz  bool
	DisableMetrics bool
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Run(ctx context.Context) error {
	addr := s.addr
	// connnect to gRPC server
	conn, err := dail("tcp", addr)
	if err != nil {
		return err
	}

	// Gateway ServeMux
	// TODO 如果使用 WithOutgoingHeaderMatcher, 需要修改下面的 handleForwardResponseServerMetadata
	var optoins []runtime.ServeMuxOption
	errHandler := s.HTTPErrorHandleFunc
	if errHandler == nil {
		errHandler = DefaultHTTPErrorHandleFunc
	}
	optoins = append(optoins, runtime.WithErrorHandler(runtimeHTTPErrorHandler(errHandler)))
	// optoins = append(optoins, runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler))

	// set jsonpb TODO 新版本细节确认
	jsonpb := &runtime.JSONPb{}
	jsonpb.MarshalOptions.UseProtoNames = true
	jsonpb.MarshalOptions.UseEnumNumbers = true
	jsonpb.MarshalOptions.EmitUnpopulated = true
	optoins = append(optoins, runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpb))
	optoins = append(optoins, runtime.WithIncomingHeaderMatcher(s.incomingHeaderMatcher()))

	gw := runtime.NewServeMux(optoins...)
	for _, reg := range s.ServerRegisters {
		if err := reg.RegisterGateway(ctx, gw, conn); err != nil {
			return err
		}
	}

	// http mux 可以把自己的 HTTP 路由添加到这里
	mux := http.NewServeMux()
	// 加载自定义的 http handler
	if fn := s.CustromHTTPHandlerSetFunc; fn != nil {
		fn(mux)
	}
	// mux.HandleFunc("/swagger/", swaggerServer(opts.SwaggerDir)) // TODO
	mux.HandleFunc("/health", healthzServer(conn))
	// pprof
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	mux.Handle("/", gw)

	var grpcOpts []grpc.ServerOption
	grpcOpts = append(grpcOpts, grpc.UnaryInterceptor(middleware.ChainUnaryServer(
		append([]grpc.UnaryServerInterceptor{
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(printPanic)),
			// trace id 拦截器
			InterceptorTraceID(),
			InterceptorLogRequest(),                           // 日志记录
			//InterceptorGRPCErrorHandle(gRPCErrorHandleFunc), // 错误处理
			//InterceptorAuth(s.lg, s.webAuth, s.secrets.GetSecret), // API鉴权
		},
			s.UnaryServerInterceptors...,
		)...,
	)))

	grpcServer := grpcUtil.NewGRPCServer(grpcOpts...)
	for _, reg := range s.ServerRegisters {
		reg.RegisterGRPC(grpcServer)
	}
	reflection.Register(grpcServer)

	h := exposeHeaders(traceId(mux))
	svr := &http.Server{
		Addr:    addr,
		Handler: grpcHandlerFunc(grpcServer, s.allowCORS(h)), // TODO wrappers, 可以变为参数
	}

	stop := make(chan error)
	go func() {
		logx.Infof(ctx, "HTTP & gRPC Server start at %q", addr)
		if err := svr.ListenAndServe(); err == nil || errors.Is(err, http.ErrServerClosed) {
			stop <- nil
		} else {
			stop <- errors.Wrap(err, "HTTP & gRPC server serve err")
		}
	}()

	select {
	case <-ctx.Done():
		logx.Warnf(ctx, "Canceled, stop HTTP server")
		cctx, cancel := context.WithTimeout(context.TODO(), s.ShutdownTimeout)
		defer cancel()
		return errors.Wrap(svr.Shutdown(cctx), "HTTP & gRPC server shutdown err")
	case err := <-stop:
		return err
	}
}

func (s *Server) incomingHeaderMatcher() func(originHttpHeaderKey string) (string, bool) {
	allow := make(map[string]bool)
	for _, v := range s.IncomingHeaderWhiteList {
		allow[strings.ToLower(v)] = true
	}
	return func(originHttpHeaderKey string) (string, bool) {
		v, ok := runtime.DefaultHeaderMatcher(originHttpHeaderKey) // 先用标注的处理
		if ok {
			return v, ok
		}
		// 处理我们自己认识的 header, 最好变成小写
		v = strings.ToLower(originHttpHeaderKey)
		if allow[v] {
			return v, true
		}
		return "", false
	}
}


func dail(network, addr string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		// grpc.WithBlock(),
	}
	switch network {
	case "tcp":
	case "unix":
		d := func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}
		opts = append(opts, grpc.WithDialer(d))
	default:
		return nil, fmt.Errorf("unsupported network type: %q", network)
	}
	return grpc.Dial(addr, opts...)
}

func gRPCErrorHandleFunc(err error) error {
	er, ok := errs.As(err)
	if !ok {
		er = errs.NewSimpleError("system.SystemException", "系统错误: "+err.Error())
	}
	return er
}

func printPanic(p interface{}) (err error) {
	var buf [4096]byte
	n := go_runtime.Stack(buf[:], false)

	msgMap := make(map[string]string)
	msgMap["level"] = "error"
	msgMap["msg"] = fmt.Sprintf("panic: %v\r\n%s", p, string(buf[:n]))
	msgJson, _ := json.Marshal(msgMap)
	fmt.Printf("%s\n", msgJson)
	return status.Errorf(codes.Internal, "系统异常")
}

// 默认的 HTTPErrorHandleFunc
func DefaultHTTPErrorHandleFunc(sts *status.Status) (_ []byte, statusCode int) {
	// 暂时统一为 500.
	statusCode = 500

	er := &errs.Error{
		Code:    "system.SystemException",
		Message: sts.Err().Error(),
	}
	if details := sts.Details(); len(details) > 0 {
		if err, ok := details[0].(*errs.Error); ok {
			er.Code = err.Code
			er.Message = err.Message
		}
	}
	b, err1 := json.Marshal(er)
	if err1 != nil {
		// 序列化错误,使用 fallback
		return []byte(`{code": "system.SystemException", "message": "系统错误: 序列化结果失败"}`), statusCode
	}
	return b, statusCode
}

// 参考 gRPC gateway 的 DeafultHTTPErrorHandler
// panic if errHandle is nil
func runtimeHTTPErrorHandler(errHandle HTTPErrorHandleFunc) runtime.ErrorHandlerFunc {
	if errHandle == nil {
		panic("HTTPErrorHandleFunc is nil")
	}

	return func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
		// return Internal when Marshal failed

		s := status.Convert(err)
		pb := s.Proto()

		w.Header().Del("Trailer")
		w.Header().Del("Transfer-Encoding")

		contentType := marshaler.ContentType(pb)
		w.Header().Set("Content-Type", contentType)

		md, ok := runtime.ServerMetadataFromContext(ctx)
		if !ok {
			grpclog.Infof("Failed to extract ServerMetadata from context")
		}

		handleForwardResponseServerMetadata(w, mux, md)

		// RFC 7230 https://tools.ietf.org/html/rfc7230#section-4.1.2
		// Unless the request includes a TE header field indicating "trailers"
		// is acceptable, as described in Section 4.3, a server SHOULD NOT
		// generate trailer fields that it believes are necessary for the user
		// agent to receive.
		var wantsTrailers bool

		if te := r.Header.Get("TE"); strings.Contains(strings.ToLower(te), "trailers") {
			wantsTrailers = true
			handleForwardResponseTrailerHeader(w, md)
			w.Header().Set("Transfer-Encoding", "chunked")
		}

		buf, st := errHandle(s)
		w.WriteHeader(st)
		if _, err := w.Write(buf); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}

		if wantsTrailers {
			handleForwardResponseTrailer(w, md)
		}
	}
}

func defaultOutgoingHeaderMatcher(key string) (string, bool) {
	return fmt.Sprintf("%s%s", runtime.MetadataHeaderPrefix, key), true
}

func handleForwardResponseServerMetadata(w http.ResponseWriter, mux *runtime.ServeMux, md runtime.ServerMetadata) {
	for k, vs := range md.HeaderMD {
		// TODO 这里是写死的! 如果 ServerMux 使用了
		if h, ok := defaultOutgoingHeaderMatcher(k); ok {
			for _, v := range vs {
				w.Header().Add(h, v)
			}
		}
	}
}

func handleForwardResponseTrailerHeader(w http.ResponseWriter, md runtime.ServerMetadata) {
	for k := range md.TrailerMD {
		tKey := textproto.CanonicalMIMEHeaderKey(fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, k))
		w.Header().Add("Trailer", tKey)
	}
}

func handleForwardResponseTrailer(w http.ResponseWriter, md runtime.ServerMetadata) {
	for k, vs := range md.TrailerMD {
		tKey := fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, k)
		for _, v := range vs {
			w.Header().Add(tKey, v)
		}
	}
}

func healthzServer(conn *grpc.ClientConn) http.HandlerFunc {
	thisIP := os.Getenv("THIS_IP")
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if s := conn.GetState(); s != connectivity.Ready && s != connectivity.Idle {
			http.Error(w, fmt.Sprintf("grpc server is %s", s), http.StatusBadGateway)
			return
		}
		fmt.Fprintf(w, "iot-eba-ok: %v\nTHIS_IP:%s\n", time.Now(), thisIP)
	}
}

func (s *Server) allowCORS(h http.Handler) http.Handler {
	allowHeaders := strings.Join(append([]string{"Content-Type", "Accept", "Authorization"}, s.IncomingHeaderWhiteList...), ",")
	//allowHeaders := strings.Join([]string{"Content-Type", "Accept", "Authorization", "token", "groupid",
	//	"app", "pf_token", "x-requested-with", trace.IotTraceHeader}, ",")
	allowMethods := strings.Join([]string{"GET", "HEAD", "POST", "PUT", "DELETE"}, ",")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
				w.Header().Set("Access-Control-Allow-Methods", allowMethods)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}