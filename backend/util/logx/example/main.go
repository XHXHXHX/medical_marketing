package main

import (
	"context"
	"github.com/RyouZhang/async-go"
	"gitlab.aiforward.cn/property_finance/common/logx"
	"gitlab.aiforward.cn/property_finance/common/logx/core"
	time_split "gitlab.aiforward.cn/property_finance/common/logx/provider/file/time"
	"gitlab.aiforward.cn/property_finance/common/logx/provider/zap"
	"gitlab.aiforward.cn/property_finance/common/recover/sentry"

	"time"
)

type A struct {
	name string
}

func standardLog() {
	ctx := context.WithValue(context.Background(), "name", "standardlog")

	logx.WithContext(ctx).AddPair("name", "hello").Error("fuck:%s|%s", "hello world", "ssa")
	logx.Debug(ctx, "standard debug message", map[string]interface{}{"field": "16"}) // 设置的日志界别为warnLevel，不打印
	logx.Info(ctx, "standard info message", map[string]interface{}{"field": []string{"key1", "key2", "key3"}})
	logx.Warn(ctx, "standard warn message", nil)
	logx.Error(ctx, "standed error message", map[string]interface{}{"field": []string{"key1", "key2", "key3"}})
	logx.WithContext(ctx)
	x := 5555
	logx.Error(ctx,
		"standard error message",
		core.NewFields().
			AddPair("field", "18").
			AddPair("flag", int64(1)).
			AddPair("duration", 123).
			AddFields("item", core.NewFields().
				AddPair("a", &A{name: "sadsad"}).
				AddPair("b", &x).
				AddPair("c", 3213)).
			Export())
}

func customLog() {
	ctx := context.WithValue(context.Background(), "name", "customlog")
	logger := logx.Get("logx_name")
	logger.Debug(ctx, "custom debug message", map[string]interface{}{"field": 25})
	logger.Info(ctx, "custom info message", map[string]interface{}{"field": 26})
	logger.Warn(ctx, "custom warn message", map[string]interface{}{"field": 27})
	logger.Error(ctx, "custom error message", nil)

	for i := 0; i < 10; i++ {
		go func(i int) {

			sl := make([]int, 0)
			for j := 0; j < i; j++ {
				sl = append(sl, j)
			}
			logger.Warn(context.Background(), "customlog warn goroutine", map[string]interface{}{"goroutine": i, "line": 37})
		}(i)

	}

}

func main() {
	//err := sentry.SetDSN("http://f0fd1dcc97f24053b49d84b9fa4518aa@sentry.com/686")
	//if err != nil {
	//	panic(err)
	//}
	// 初始化logx namespace，project，全局设置一次
	logx.InitNamespaceProject("pf", "subject")

	// 设置logx panicHandler 全局设置一次
	async.SetPanicHandler(sentry.Handler)

	// 使用logx默认的logger
	logx.SetLevel(core.WarnLevel) // 设置打印日志界别

	// 添加logrus provider,logger之间共享provider
	logx.AddProvider(zap.NewProvider())
	logx.AddProvider(time_split.NewProvider("/tmp", "time_split", ""))

	// logx默认logger增加filter，filter之间不共享
	logx.AddFilter(func(ctx context.Context) map[string]interface{} { // 注册context filter
		name := ctx.Value("name")
		age := ctx.Value("age")
		m := map[string]interface{}{
			"name": name,
			"age":  age,
		}
		// panic(m)
		return m
	})

	// 测试截断
	longString := "默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level" +
		"默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level默认使用logx.SetLevel()设置的日志界别 为logger设置level"
	logx.Warn(context.Background(), longString, nil)
	logx.Warn(context.Background(), "test", map[string]interface{}{
		"long_string": longString,
	})
	// 通用打日志方法
	// standardLog()
	// 使用自定义logger，logger之间不共享filter
	//logger := logx.Get("logx_name") // 传入name获取logger
	//logger.SetLevel(core.InfoLevel) // 默认使用logx.SetLevel()设置的日志界别 为logger设置level
	//logger.AddFilter(func(ctx context.Context) map[string]interface{} { // 注册context filter
	//	name := ctx.Value("name")
	//	m := map[string]interface{}{
	//		"name": name,
	//	}
	//	return m
	//})
	// 自定义打日志方法
	// customLog()
	time.Sleep(10 * time.Second)
}
