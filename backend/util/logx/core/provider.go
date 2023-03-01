package core

import (
	"context"
	"sync"

	"github.com/RyouZhang/async-go"
)

var (
	providers    []*Provider
	pLokcer      sync.Mutex
	syncProvider bool = false
)

func init() {
	providers = make([]*Provider, 0)
}

type Provider struct {
	process func(int64, Level, string, map[string]interface{})
	input   chan *logData
}

func (p *Provider) run() {
	for data := range p.input {
		_, err := async.Safety(func() (interface{}, error) {
			args := map[string]interface{}{}
			for k, v := range data.args {
				args[k] = v
			}
			args["project"] = data.project
			args["namespace"] = data.namespace
			args["logger"] = data.logger
			args["linenum"] = data.linenum
			p.process(data.ts, data.level, data.msg, args)
			return nil, nil
		})
		if err != nil {
			panicLogger.Error(context.Background(), err.Error(), nil)
		}
	}
}

func NewProvider(method func(int64, Level, string, map[string]interface{})) *Provider {
	return &Provider{
		process: method,
		input:   make(chan *logData, 16),
	}
}

func AddProvider(p *Provider) {
	pLokcer.Lock()
	defer pLokcer.Unlock()
	providers = append(providers, p)
	go p.run()
}

func sendProvider(data *logData) {
	pLokcer.Lock()
	defer pLokcer.Unlock()

	for _, p := range providers {
		if syncProvider {
			p.input <- data
		} else {
			select {
			case p.input <- data:
			default: //do nothing
			}
		}
	}
}

func SetProviderSync() {
	syncProvider = true
}
