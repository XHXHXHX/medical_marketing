package server

import (
	"context"
	"time"

	"github.com/XHXHXHX/medical_marketing/util/common"
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
//func InterceptorGRPCErrorHandle(handle GRPCErrorHandleFunc) grpc.UnaryServerInterceptor {
//	if handle == nil {
//		panic("handle func is nil")
//	}
//
//	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
//		resp, err := handler(ctx, req)
//		if err == nil {
//			return resp, nil
//		}
//		sts, _ := status.FromError(err)
//		if sts == nil {
//			sts = status.New(codes.Internal, "Custom error")
//		}
//		detail := handle(err)
//		sts1, err1 := sts.WithDetails(detail)
//		if err1 != nil {
//			sts1 = sts
//		}
//		return nil, sts1.Err()
//	}
//}

// 权限认证
//func InterceptorAuth(logger *zap.Logger, auth authority.WebAuthService, getAPISecret func(string) (string, error)) grpc.UnaryServerInterceptor {
//	if logger == nil {
//		panic("logger is nil")
//	}
//
//	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
//		// 心跳检测，跳过权限认证
//		if strings.HasPrefix(info.FullMethod, "/iot_eba.v1api.calc_model.V1APICalcModel/HeartBeat") {
//			return handler(ctx, req)
//		}
//		// TODO 临时测试
//		if strings.HasPrefix(info.FullMethod, "/iot_eba.v1api.mobile.V1APIMobile") {
//			return handler(ctx, req)
//		}
//		// wapi
//		if strings.HasPrefix(info.FullMethod, "/iot_eba.wapi.") {
//			checker := func(ctx context.Context, getAPISecret authutil.SecretProvider) error {
//				params, err := authutil.GetAuthParams(ctx)
//				if err != nil {
//					return err
//				}
//				secret, err := getAPISecret(params.AppKey)
//				if err != nil {
//					return err
//				}
//				if err := params.Validate(secret); err != nil {
//					return err
//				}
//
//				if m := meta.GetMetaFromContext(ctx); m != nil {
//					m.AddLogFields(logutil.AppKey(params.AppKey))
//				}
//
//				return nil
//			}
//
//			if checker(ctx, getAPISecret) != nil {
//				return nil, status.Errorf(codes.Unauthenticated, "auth failed")
//			}
//			return handler(ctx, req)
//		}
//		startTime := time.Now()
//
//		// 权限认证
//		// 假设,我们使用的 HTTP header 是 token, 映射到 grpc 是
//		// 取出 header
//
//		user, err := auth.GetUserInfo(ctx)
//		if err != nil {
//			return nil, errs.AuthTokenError
//		}
//
//		if m := meta.GetMetaFromContext(ctx); m != nil {
//			m.AddLogFields(logutil.JSONAny("user-info", user))
//		}
//
//		logger.Info("get token result", zap.Int("user_id", user.UserID), logutil.TakesMSFrom(startTime))
//		ctx = auth.SetUserInfo(ctx, user)
//		// 操作权限认证, url和操作权限映射
//
//		if _, ok := skipMethods[info.FullMethod]; ok {
//			return handler(ctx, req)
//		}
//		code, ok := methodActionMap[info.FullMethod]
//		if !ok {
//			return nil, errs.NoAuthority
//		}
//
//		if err := auth.CheckCASPermission(ctx, 0, code, user.UserID); err != nil {
//			logger.Warn("user has no permission",
//				logutil.JSONAny("user-info", user),
//				zap.String("err", err.Error()),
//			)
//			return nil, errs.NoAuthority.WithStack()
//		}
//
//		return handler(ctx, req)
//	}
//}

// traceID
func InterceptorTraceID() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx = common.SetGlobalID(ctx)

		return handler(ctx, req)
	}
}

