package logx

import (
	"context"
	"fmt"
	"sync"

	"github.com/XHXHXHX/medical_marketing/util/logx/core"
)

var (
	loggers   sync.Map
	stdOnce   sync.Once
	stdLogger *core.Logger
)

func InitNamespaceProject(namespace string, project string) {
	core.Namespace, core.Project = namespace, project
}

func getStdLogger() *core.Logger {
	stdOnce.Do(func() {
		stdLogger = core.NewLogger("__default__")
	})
	return stdLogger
}

func GetDefault() *core.Logger {
	return getStdLogger()
}

// Debug _
func Debug(ctx context.Context, msg string, args map[string]interface{}) {
	getStdLogger().Debug(CtxWithSkipOffset(ctx), msg, args)
}

// Info _
func Info(ctx context.Context, msg string, args map[string]interface{}) {
	getStdLogger().Info(CtxWithSkipOffset(ctx), msg, args)
}

// Info _
func Infof(ctx context.Context, msg string, args ...interface{}) {
	getStdLogger().Info(CtxWithSkipOffset(ctx), fmt.Sprintf(msg, args...), nil)
}

// Warn _
func Warn(ctx context.Context, msg string, args map[string]interface{}) {
	getStdLogger().Warn(CtxWithSkipOffset(ctx), msg, args)
}

// Warn _
func Warnf(ctx context.Context, msg string, args ...interface{}) {
	getStdLogger().Warn(CtxWithSkipOffset(ctx), fmt.Sprintf(msg, args...), nil)
}

// Error _
func Error(ctx context.Context, msg string, args map[string]interface{}) {
	getStdLogger().Error(CtxWithSkipOffset(ctx), msg, args)
}

// Error _
func Errorf(ctx context.Context, msg string, args ...interface{}) {
	getStdLogger().Error(CtxWithSkipOffset(ctx), fmt.Sprintf(msg, args...), nil)
}

func CtxWithSkipOffset(ctx context.Context) context.Context {
	old := ctx.Value(core.SkipOffsetKey)
	if old != nil {
		ctx = context.WithValue(ctx, core.SkipOffsetKey, old.(int)+1)
	} else {
		ctx = context.WithValue(ctx, core.SkipOffsetKey, 1)
	}
	return ctx
}

func Get(name string) *core.Logger {
	target, ok := loggers.Load(name)
	if ok {
		return target.(*core.Logger)
	}
	logger := core.NewLogger(name)
	loggers.Store(logger.Name, logger)
	return logger
}

func SetLevel(level core.Level) {
	core.LogLevel = level
	if stdLogger != nil {
		stdLogger.SetLevel(level)
	}
}

func SetLevelName(level string) {
	switch level {
	case "info":
		SetLevel(core.InfoLevel)
	case "warn":
		SetLevel(core.WarnLevel)
	case "error":
		SetLevel(core.ErrorLevel)
	case "debug":
		SetLevel(core.DebugLevel)
	default:
		SetLevel(core.InfoLevel)
	}
}

func SetSyncMode() {
	core.SetProviderSync()
}

func AddHook(hook core.Hook) {
	getStdLogger().AddHook(hook)
}

func AddProvider(p *core.Provider) {
	core.AddProvider(p)
}

func AddFilter(filter func(context.Context) map[string]interface{}) {
	getStdLogger().AddFilter(filter)
}

func SetFormater(f core.Formater) {
	getStdLogger().SetFormater(f)
}

func RegisterRule(key string, vType core.FieldType) error {
	return core.RegisterRule(key, vType)
}

// with new log action
func WithContext(ctx context.Context) *core.Action {
	return getStdLogger().WithContext(CtxWithSkipOffset(ctx))
}
