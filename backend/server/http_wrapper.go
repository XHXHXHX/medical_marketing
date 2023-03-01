package server

import (
	"net/http"
	"sort"
	"strings"

	"github.com/XHXHXHX/medical_marketing/util/trace"
)

type interceptor = func(http.Handler) http.Handler

// 对 http server 的包装，封装了 拦截器
type server struct {
	httpMux      *http.ServeMux
	interceptors []interceptor
}

// 初始化，传入定义的拦截器
func NewWithInterceptors(interceptors ...interceptor) *server {
	// 逆序，排前边的拦截器先执行
	sort.Slice(interceptors, func(i, j int) bool {
		return j < i
	})
	return &server{
		httpMux:      http.NewServeMux(),
		interceptors: interceptors,
	}
}

// 添加 http 路由，配置拦截器
func (s *server) AddRoute(path string, handler http.Handler) {
	h := handler
	for _, c := range s.interceptors {
		h = c(h)
	}
	s.httpMux.Handle(path, h)
}

// 添加 http 路由，配置自定义的拦截器
func (s *server) AddRouteWith(path string, handler http.Handler, interceptors ...interceptor) {
	sort.Slice(interceptors, func(i, j int) bool {
		return j < i
	})
	h := handler
	for _, c := range interceptors {
		h = c(h)
	}
	s.httpMux.Handle(path, h)
}

// 实现 http.Handler 接口
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.httpMux.ServeHTTP(w, r)
}

// 新增拦截器：添加 trace id header
func traceId(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// 设置traceid
		traceID := trace.GetTraceFromHTTP(r)
		if len(traceID) == 0 {
			traceID = trace.NewTraceID()
		}
		trace.SetHTTPResponseTraceID(w, traceID)

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// 新增拦截器：允许前端获取的 header
func exposeHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 允许前端获取 traceid
		traceHeader := upperTraceHeader()
		exposeHeaders := strings.Join([]string{traceHeader, "Grpc-Metadata-" + traceHeader}, ",")
		w.Header().Set("Access-Control-Expose-Headers", exposeHeaders)
		h.ServeHTTP(w, r)
	})
}

func upperTraceHeader() string {
	return firstUpper(trace.IotTraceHeader)
}

func firstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
