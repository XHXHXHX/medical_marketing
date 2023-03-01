package grpc

import (
	"context"
	"crypto/tls"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"math"
	"time"

	"google.golang.org/grpc/metadata"
)

const (
	grpcMaxRecvMsgSize = 10 << 20
	maxStreams         = math.MaxUint32
	maxSendBytes       = math.MaxInt32
)

func HeaderFromContext(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	return getHeader(md, key)
}

func getHeader(md metadata.MD, key string) string {
	if v := md.Get("grpcgateway-" + key); len(v) > 0 {
		return v[0]
	}
	// 自定义header 没有前缀
	if v := md.Get(key); len(v) > 0 {
		return v[0]
	}
	return ""
}


// Dail 默认是 insecure 的 plaintext 通信
func Dail(ctx context.Context, addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts,
		grpc.WithInsecure(),
		// grpc.WithBlock(), // 默认不要 block
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxRecvMsgSize)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             time.Second,
			PermitWithoutStream: true,
		}),
	)

	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// DialConfig config of dial
type DialConfig struct {
	// 地址
	Addr string
	// 是否明文(非tsl)
	PlainText bool

	// tsl相关配置
	// tsl 是否跳过验证
	Insecure bool
	// tsl 覆盖 server name
	OverrideServerName string

	WithBlock bool

	// 其他自定义的 options
	ExtraOptions []grpc.DialOption
}

// DialWith 类似 Dial, 支持常用的参数, 包括明文和 tsl
func DialWith(ctx context.Context, cfg DialConfig) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxRecvMsgSize)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	)

	if len(cfg.Addr) == 0 {
		return nil, errors.New("grpc addr is empty")
	}

	if cfg.PlainText {
		opts = append(opts, grpc.WithInsecure())
	} else {
		var tlsConf tls.Config
		if cfg.Insecure {
			tlsConf.InsecureSkipVerify = true
		}
		creds := credentials.NewTLS(&tlsConf)
		if len(cfg.OverrideServerName) != 0 {
			creds.OverrideServerName(cfg.OverrideServerName)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	if cfg.WithBlock {
		opts = append(opts, grpc.WithBlock())
	}

	opts = append(opts, cfg.ExtraOptions...)

	conn, err := grpc.DialContext(ctx, cfg.Addr, opts...)
	if err != nil {
		return nil, err
	}
	return conn, err
}

// 常用的 server 创建
func NewGRPCServer(options ...grpc.ServerOption) *grpc.Server {
	var grpcOpts = []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    10 * time.Second,
			Timeout: 5 * time.Second,
		}),
	}
	grpcOpts = append(grpcOpts, options...)
	return grpc.NewServer(grpcOpts...)
}
