package zap

import (
	"os"

	"github.com/XHXHXHX/medical_marketing/util/logx/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func makeFields(args map[string]interface{}) []zapcore.Field {
	result := make([]zapcore.Field, 0)
	for key, value := range args {
		result = append(result, zap.Any(key, value))
	}
	return result
}

func NewProvider() *core.Provider {
	ec := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "linenum",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	zc := zapcore.NewCore(
		zapcore.NewJSONEncoder(ec),
		zapcore.AddSync(os.Stdout),
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.DebugLevel
		}),
	)

	logger := zap.New(zc)
	defer logger.Sync()
	return core.NewProvider(func(ts int64, level core.Level, msg string, args map[string]interface{}) {
		fields := makeFields(args)
		switch level {
		case core.DebugLevel:
			{
				logger.Debug(msg, fields...)
			}
		case core.InfoLevel:
			{
				logger.Info(msg, fields...)
			}
		case core.WarnLevel:
			{
				logger.Warn(msg, fields...)
			}
		case core.ErrorLevel:
			{
				logger.Error(msg, fields...)
			}
		default:
			{
				//do nothing
			}
		}
	})
}
