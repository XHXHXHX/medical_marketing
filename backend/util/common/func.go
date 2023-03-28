package common

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/metadata"
)

func SetUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, USER_ID, userID)
}

func GetUserID(ctx context.Context) int64 {
	if v := ctx.Value(USER_ID); v != nil {
		return v.(int64)
	}
	return 0
}

func SetAdmin(ctx context.Context, isAdmin bool) context.Context {
	return context.WithValue(ctx, IS_ADMIN, isAdmin)
}

func GetAdmin(ctx context.Context) bool {
	if v := ctx.Value(IS_ADMIN); v != nil {
		return v.(bool)
	}
	return false
}

func GetUserName(ctx context.Context) string {
	if v := ctx.Value(USER_NAME); v != nil {
		return v.(string)
	}
	return ""
}

func SetRole(ctx context.Context, role int64) context.Context {
	return context.WithValue(ctx, ROLE, role)
}

func GetRole(ctx context.Context) int64 {
	if v := ctx.Value(ROLE); v != nil {
		return v.(int64)
	}
	return 0
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

func MD5Password(passwd, salt string) string {
	return MD5(fmt.Sprintf("%s%s", passwd, salt))
}

func MD5(s string) string {
	m := md5.New()
	m.Write([]byte(s))
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

// ascii 随机
func RandString(length int) string {
	dice := 126 - 32

	var list []string

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		r := rand.Intn(dice) + 32
		list = append(list, string(rune(r)))
	}

	return strings.Join(list, "")
}

// ascii 随机
func RandNormalString(length int) string {
	dice := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var list []rune
	l := len(dice)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		r := rand.Intn(l-1)

		list = append(list, rune(dice[r]))
	}

	return string(list)
}


// 类型转换
func InterfaceToString(s interface{}) string {
	switch v := s.(type) {
	case string:
		return v
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case []string, []int32, []int64, []float64, []float32, map[string]int32, map[string]int64, map[string]string, map[string]float32, map[string]float64, []map[string]int32, []map[string]int64, []map[string]string, []map[string]float32, []map[string]float64:
		jsonData, _ := json.Marshal(v)
		return string(jsonData)
	default:
		jsonData, _ := json.Marshal(v)
		return string(jsonData)
	}
}

func InterfaceToInt(s interface{}) (int64, error) {
	switch v := s.(type) {
	case string:
		return strconv.ParseInt(v, 10, 64)
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	default:
		return 0, errors.New("unrecognized type")
	}
}