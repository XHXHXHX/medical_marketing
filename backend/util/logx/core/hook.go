package core

import (
	"context"

	"github.com/RyouZhang/async-go"
)

// Hook _
type Hook interface {
	Check(l Level) bool
	Process(int64, Level, string, map[string]interface{}) error
}

type hookHost struct {
	h     Hook
	input chan *logData
}

func newHookHost(h Hook) *hookHost {
	return &hookHost{
		h:     h,
		input: make(chan *logData, 32),
	}
}

func (h *hookHost) run() {
	for data := range h.input {

		_, err := async.Safety(func() (interface{}, error) {
			pErr := h.h.Process(data.ts, data.level, data.msg, data.args)
			return nil, pErr
		})
		if err != nil {
			panicLogger.Error(context.Background(), err.Error(), nil)
		}
	}
}
