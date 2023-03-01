package common

import (
	"context"
	"crypto/md5"
	"encoding/hex"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/metadata"
)

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, USER_ID, userID)
}

func GetUserID(ctx context.Context) string {
	if v := ctx.Value(USER_ID); v != nil {
		return v.(string)
	}
	return ""
}

func GetUserName(ctx context.Context) string {
	if v := ctx.Value(USER_NAME); v != nil {
		return v.(string)
	}
	return ""
}

func SetUserName(ctx context.Context, userName string) context.Context {
	return context.WithValue(ctx, USER_NAME, userName)
}

func GetGlobalID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		logID := GetStrValueFromMD(md, LOG_ID)
		if logID != "" {
			return logID
		}
	}
	if v := ctx.Value(LOG_ID); v != nil {
		return v.(string)
	}
	return ""
}

func GenerateGlobalID() string {
	value := uuid.Must(uuid.NewV4(), nil).String()
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

func SetGlobalID(ctx context.Context) context.Context {
	return context.WithValue(ctx, LOG_ID, GenerateGlobalID())
}

func GetStrValueFromMD(md metadata.MD, key string) string {
	if v := md.Get(key); len(v) > 0 {
		return v[0]
	}
	return ""
}