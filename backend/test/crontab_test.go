package tests

import (
	"context"
	"github.com/XHXHXHX/medical_marketing/config"
	"github.com/XHXHXHX/medical_marketing/service/report"
	"github.com/XHXHXHX/medical_marketing/util/logx/core"
	"go.uber.org/dig"
	"testing"
)

var env = "dev.yaml"

type Startup struct {
	dig.In

	// 这里添加你需要使用的类型,会自动初始化
	LoggerX *core.Logger
	Report report.Service
}

func setup() *Startup {
	cfg := GetConfig(env)
	var p Startup
	if err := config.NewContainer(&cfg).Invoke(func(s Startup) {
		p = s
	}); err != nil {
		panic(err)
	}

	return &p
}

func GetConfig(env string) (cfg config.Config) {
	err := new(Builder).FromYaml(env).WithENV().Build(&cfg)
	if err != nil {
		panic(err)
	}
	return
}

func TestAutoChangeBelong(t *testing.T) {
	svr := setup()

	ctx := context.Background()

	svr.Report.AutoChangeBelong(ctx)
}