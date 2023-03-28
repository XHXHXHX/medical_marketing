package config

import (
	"context"
	"github.com/XHXHXHX/medical_marketing/service/customer_task"
	"time"

	customerTaskRepo "github.com/XHXHXHX/medical_marketing/repository/customer_task"
	customerTaskHistoryRepo "github.com/XHXHXHX/medical_marketing/repository/customer_task_history"
	reportRepo "github.com/XHXHXHX/medical_marketing/repository/report"
	userRepo "github.com/XHXHXHX/medical_marketing/repository/user"
	"github.com/XHXHXHX/medical_marketing/server"
	"github.com/XHXHXHX/medical_marketing/service/report"
	customerTaskImpl "github.com/XHXHXHX/medical_marketing/service/report/impl"
	reportImpl "github.com/XHXHXHX/medical_marketing/service/customer_task/impl"
	"github.com/XHXHXHX/medical_marketing/service/user"
	userImpl "github.com/XHXHXHX/medical_marketing/service/user/impl"
	"github.com/XHXHXHX/medical_marketing/util/common"
	"github.com/XHXHXHX/medical_marketing/util/logx"
	"github.com/XHXHXHX/medical_marketing/util/logx/core"
	"github.com/XHXHXHX/medical_marketing/util/logx/provider/zap"
	"github.com/XHXHXHX/medical_marketing/util/mysql"

	goredis "github.com/go-redis/redis/v8"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func NewContainer(cfg *Config) *dig.Container {
	container := dig.New()

	mustProvide := func(constructor interface{}, opts ...dig.ProvideOption) {
		if err := container.Provide(constructor, opts...); err != nil {
			panic(err)
		}
	}
	mustProvide(func() *Config {
		return cfg
	})

	// 注册
	mustProvide(NewMysql)
	mustProvide(NewNewLogger)
	mustProvide(NewRedisClient)

	mustProvide(NewServer)
	mustProvide(server.NewCommon)

	mustProvide(customerTaskRepo.NewRepo)
	mustProvide(customerTaskHistoryRepo.NewRepo)
	mustProvide(customerTaskImpl.NewService)

	mustProvide(userRepo.NewRepo)
	mustProvide(userImpl.NewService)

	mustProvide(reportRepo.NewRepo)
	mustProvide(reportImpl.NewService)

	return container
}

func NewServer(p struct {
	dig.In

	Cfg *Config
	Report report.Service
	User user.Service
	CustomerTask customer_task.Service
}) *server.Server {
	return server.NewServer(p.Cfg.Server.Addr, p.Report, p.User, p.CustomerTask)
}

func NewMysql(cfg *Config) (*gorm.DB, error) {
	return mysql.NewMysql(
		cfg.Mysql.Host,
		cfg.Mysql.Username,
		cfg.Mysql.Password,
		cfg.Mysql.DBName,
		"",
	)
}

// NewNewLogger 创建 Logger 实例
func NewNewLogger() (*core.Logger, error) {
	logx.InitNamespaceProject("backend", "sale")
	logx.SetLevelName("info")
	logx.AddProvider(zap.NewProvider())
	// 注册context filter
	logx.AddFilter(func(ctx context.Context) map[string]interface{} {
		return map[string]interface{}{
			common.USER_NAME:  common.GetUserName(ctx),
			common.USER_ID:    common.GetUserID(ctx),
			common.LOG_ID:     common.GetGlobalID(ctx),
		}
	})
	return logx.GetDefault(), nil
}

func NewRedisClient(cfg *Config) (*goredis.Client, error) {
	c := cfg.Redis
	logx.Infof(context.Background(), "Redis config host:[%s], pass(len):%d",
		c.Host, len(c.Passwd))
	client := goredis.NewClient(&goredis.Options{
		Addr:     c.Host,
		Password: c.Passwd,
		DB:       c.DB,
	})
	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}