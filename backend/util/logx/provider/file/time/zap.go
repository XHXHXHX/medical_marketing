package time_split

import (
	"gitlab.aiforward.cn/property_finance/common/logx/core"

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

//example dateformat "2006-01-02T15-04-05.000"
//daily	split	20060102
//hour split	2006010215
//minute split	200601021504
func NewProvider(basePath string, prefix string, dateFormat string) *core.Provider {
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

	if len(dateFormat) == 0 {
		dateFormat = "20060102"
	}

	zc := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(ec),
			zapcore.AddSync(&splitFileLogger{
				BaseDir:    basePath,
				Prefix:     prefix + ".log.info",
				DateFormat: dateFormat,
			}),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.DebugLevel
			})),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(ec),
			zapcore.AddSync(&splitFileLogger{
				BaseDir:    basePath,
				Prefix:     prefix + ".log.error",
				DateFormat: dateFormat,
			}),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.ErrorLevel
			})),
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
