package config

import (
	"context"

	reportRepo "github.com/XHXHXHX/medical_marketing/repository/report"
	"github.com/XHXHXHX/medical_marketing/server"
	reportService "github.com/XHXHXHX/medical_marketing/service/report"
	reportImpl "github.com/XHXHXHX/medical_marketing/service/report/impl"
	"github.com/XHXHXHX/medical_marketing/util/common"
	"github.com/XHXHXHX/medical_marketing/util/logx"
	"github.com/XHXHXHX/medical_marketing/util/logx/core"
	"github.com/XHXHXHX/medical_marketing/util/logx/provider/zap"
	"github.com/XHXHXHX/medical_marketing/util/mysql"

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

	mustProvide(NewServer)
	mustProvide(NewCommonServer)

	mustProvide(NewReportRepository)
	mustProvide(NewReportService)
	mustProvide(NewReportServer)

	mustProvide(server.NewReportServer, dig.Group("server_registers"))

	return container
}

func NewServer(cfg *Config) *server.Server {
	return server.NewServer(cfg.Server.Addr)
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

func NewCommonServer(report reportService.Service) *server.Common {
	return server.NewCommon(report)
}

func NewReportRepository(client *gorm.DB) reportRepo.Repository {
	return reportRepo.NewRepo(client)
}

func NewReportService(repo reportRepo.Repository) reportService.Service {
	return reportImpl.NewService(repo)
}

func NewReportServer(c *server.Common) *server.ReportServer {
	return server.NewReportServer(c)
}