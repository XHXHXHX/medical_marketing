package main

import (
	"flag"
	"fmt"
	"context"
	"github.com/XHXHXHX/medical_marketing/util/logx/core"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/XHXHXHX/medical_marketing/util/logx"
	"github.com/XHXHXHX/medical_marketing/config"
	"github.com/XHXHXHX/medical_marketing/server"
	"github.com/XHXHXHX/medical_marketing/util/conf"


	"go.uber.org/dig"
	"go.uber.org/zap"
)

type Args struct {
	ConfigFile string
}

func main() {
	var args Args
	flag.StringVar(&args.ConfigFile, "config", "dev.yaml", "Specify config file name")
	flag.Parse()
	if v := os.Getenv("CONFIG_FILE"); len(v) != 0 {
		args.ConfigFile = v
	}

	fmt.Println("using config:", args.ConfigFile)

	var cfg config.Config
	err := new(conf.Builder).FromYaml(args.ConfigFile).WithENV().Build(&cfg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", cfg)

	type Startup struct {
		dig.In

		// 这里添加你需要使用的类型,会自动初始化
		Server               *server.Server
		LoggerX *core.Logger
	}
	var p Startup
	if err := config.NewContainer(&cfg).Invoke(func(s Startup) {
		p = s
	}); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	wg := wait(
		func() {
			defer cancel()
			if err := p.Server.Run(ctx); err != nil {
				logx.Warnf(ctx, "server run err", zap.String("err", err.Error()))
			}
		},
	)

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
	case sg := <-quit:
		logx.Infof(ctx, "Receive signal %v and shutdown...", sg)
		cancel()
	}

	wg.Wait()
	logx.Infof(ctx, "All shutdown, exit")
}

func wait(funcs ...func()) *sync.WaitGroup {
	wg := new(sync.WaitGroup)

	wg.Add(len(funcs))
	for _, fn := range funcs {
		fn := fn
		go func() {
			defer wg.Done()
			fn()
		}()
	}

	return wg
}