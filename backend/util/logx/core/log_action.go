package core

import (
	"context"
	"fmt"
)

type Action struct {
	logger *Logger
	ctx    context.Context
	args   *Fields
	skip   int
}

func newAction(ctx context.Context, logger *Logger) *Action {
	return &Action{ctx: ctx, logger: logger, args: NewFields()}
}

func (a *Action) AddPair(key string, value interface{}) *Action {
	a.args.AddPair(key, value)
	return a
}

func (a *Action) AddArray(key string, values ...interface{}) *Action {
	a.args.AddArray(key, values...)
	return a
}

func (a *Action) AddFields(key string, value *Fields) *Action {
	a.args.AddFields(key, value)
	return a
}

func (a *Action) Debug(msg string, args ...interface{}) {
	a.logger.Debug(a.ctx, fmt.Sprintf(msg, args...), a.args.Export())
}

func (a *Action) Info(msg string, args ...interface{}) {
	a.logger.Info(a.ctx, fmt.Sprintf(msg, args...), a.args.Export())
}

func (a *Action) Warn(msg string, args ...interface{}) {
	a.logger.Warn(a.ctx, fmt.Sprintf(msg, args...), a.args.Export())
}

func (a *Action) Error(msg string, args ...interface{}) {
	a.logger.Error(a.ctx, fmt.Sprintf(msg, args...), a.args.Export())
}
