package core

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/RyouZhang/async-go"
)

type Logger struct {
	Name     string
	level    Level
	lock     sync.RWMutex
	hooks    []*hookHost
	filters  []func(context.Context) map[string]interface{}
	formater Formater
}

var (
	panicLogger = NewLogger("__logx_panic__")
)

func NewLogger(name string) *Logger {
	logger := &Logger{
		Name:     name,
		level:    LogLevel,
		filters:  make([]func(context.Context) map[string]interface{}, 0),
		hooks:    make([]*hookHost, 0),
		formater: &RuleFormater{},
	}
	return logger
}

func (l *Logger) processFilter(ctx context.Context) map[string]interface{} {
	l.lock.RLock()
	defer l.lock.RUnlock()

	result := make(map[string]interface{})
	for _, filter := range l.filters {
		temp, err := async.Safety(func() (interface{}, error) {
			r := filter(ctx)
			return r, nil
		})
		if err != nil {
			panicLogger.Error(context.Background(), err.Error(), nil)
			continue
		}
		ret, ok := temp.(map[string]interface{})
		if !ok {
			panicLogger.Warn(context.Background(), fmt.Sprintf("error type %T", temp), nil)
			continue
		}
		for key, value := range ret {
			result[key] = value
		}
	}
	return result
}

func (l *Logger) log(ctx context.Context, level Level, msg string, args map[string]interface{}) {
	lineNum := ""
	result := l.processFilter(ctx)
	if args != nil {
		for key, value := range args {
			if key == "linenum" {
				lineNum = fmt.Sprintf("%s", value)
			} else {
				result[key] = value
			}
		}
	}
	if lineNum == "" {
		skip := skipOffset
		if temp := ctx.Value(SkipOffsetKey); temp != nil {
			num, ok := temp.(int)
			if ok {
				skip += num
			}
		}
		_, file, line, ok := runtime.Caller(skip)
		if ok {
			lineNum = fmt.Sprintf("%s:%d", file, line)
		}
	}
	ts := time.Now().Unix()
	if len(l.hooks) > 0 {
		hookData := &logData{
			namespace: Namespace,
			project:   Project,
			logger:    l.Name,
			linenum:   lineNum,
			ts:        ts,
			level:     level,
			msg:       msg,
			args:      result,
		}
		l.sendHook(level, hookData)
	}

	if len(msg) > msgSize {
		msg = msg[:msgSize] + "..."
	}

	if l.formater != nil {
		result = l.formater.Format(result)
	}

	data := &logData{
		namespace: Namespace,
		project:   Project,
		logger:    l.Name,
		linenum:   lineNum,
		ts:        ts,
		level:     level,
		msg:       msg,
		args:      result,
	}
	sendProvider(data)
}

func (l *Logger) sendHook(level Level, data *logData) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for _, hook := range l.hooks {
		if hook.h.Check(level) {
			select {
			case hook.input <- data:
				{

				}
			default:
				{
					//do nothing
				}
			}
		}
	}
}

func (l *Logger) Debug(ctx context.Context, msg string, args map[string]interface{}) {
	if l.level > DebugLevel {
		return
	}
	l.log(ctx, DebugLevel, msg, args)
}

func (l *Logger) Info(ctx context.Context, msg string, args map[string]interface{}) {
	if l.level > InfoLevel {
		return
	}
	l.log(ctx, InfoLevel, msg, args)
}

func (l *Logger) Warn(ctx context.Context, msg string, args map[string]interface{}) {
	if l.level > WarnLevel {
		return
	}
	l.log(ctx, WarnLevel, msg, args)
}

func (l *Logger) Error(ctx context.Context, msg string, args map[string]interface{}) {
	if l.level > ErrorLevel {
		return
	}
	l.log(ctx, ErrorLevel, msg, args)
}

func (l *Logger) AddFilter(filter func(context.Context) map[string]interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.filters = append(l.filters, filter)
}

func (l *Logger) AddHook(h Hook) {
	l.lock.Lock()
	defer l.lock.Unlock()
	ins := newHookHost(h)
	l.hooks = append(l.hooks, ins)
	go ins.run()
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) SetFormater(f Formater) {
	l.formater = f
}

//support new log action
func (l *Logger) WithContext(ctx context.Context) *Action {
	return newAction(ctx, l)
}
