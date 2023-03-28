package server

import (
	"context"
	"time"

	"github.com/XHXHXHX/medical_marketing/errs"
	"github.com/XHXHXHX/medical_marketing/util/common"
	grpcutil "github.com/XHXHXHX/medical_marketing/util/grpc"
	"github.com/XHXHXHX/medical_marketing/util/logx"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InterceptorLogRequest 记录请求
func InterceptorLogRequest() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		starttime := time.Now()
		var (
			resp interface{}
			err  error
		)
		defer func() {
			takes := time.Since(starttime)

			code := codes.OK
			if err != nil {
				if sts, ok := status.FromError(err); ok {
					code = sts.Code()
				} else {
					code = codes.Internal
				}
			}
			if err == nil {
				logx.Info(ctx, "grpc request log OK",
					map[string]interface{}{
						"full-method": info.FullMethod,
						"code-desc": code.String(),
						"cost-ms": int(takes.Seconds()*1000),
						"req": req,
						"resp": resp,
					})
				return
			}
			logx.Warn(ctx, "grpc request log ERROR", map[string]interface{}{
				"full-method": info.FullMethod,
				"code-desc": code.String(),
				"cost-ms": int(takes.Seconds()*1000),
				"req": req,
			})
		}()

		resp, err = handler(ctx, req)
		return resp, err
	}
}

// 错误转换
func InterceptorGRPCErrorHandle(handle GRPCErrorHandleFunc) grpc.UnaryServerInterceptor {
	if handle == nil {
		panic("handle func is nil")
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}
		sts, _ := status.FromError(err)
		if sts == nil {
			sts = status.New(codes.Internal, "Custom error")
		}
		detail := handle(err)

		sts1, err1 := sts.WithDetails(detail)
		if err1 != nil {
			sts1 = sts
		}

		return nil, sts1.Err()
	}
}

// 权限认证
func InterceptorAuth(auth Auth) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 跳过认证
		if _, ok := SkipInterfaceMap[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		token := grpcutil.HeaderFromContext(ctx, common.Token)

		user, err := auth(ctx, token)
		if err != nil {
			return nil, errs.ExpireToken
		}

		ctx = common.SetUserID(ctx, user.ID)
		ctx = common.SetUserName(ctx, user.Name)
		ctx = common.SetRole(ctx, int64(user.Role))
		ctx = common.SetAdmin(ctx, user.IsAdmin())

		// 超管
		if user.IsAdmin() {
			return handler(ctx, req)
		}

		// 特定操作要求指定权限
		if roleMap, ok := RoleInterfaceMap[info.FullMethod]; ok {
			if _, ok := roleMap[user.Role]; !ok {
				return nil, errs.NoRole
			}
		}

		return handler(ctx, req)
	}
}

// traceID
func InterceptorTraceID() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx = common.SetGlobalID(ctx)

		return handler(ctx, req)
	}
}

