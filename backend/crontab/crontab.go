package crontab

import (
	"fmt"
	"runtime"

	"github.com/XHXHXHX/medical_marketing/service/report"
	"github.com/robfig/cron"
)

type Crontab struct {
	report report.Service
}

func NewCrontab( report report.Service) (*Crontab, error) {
	return &Crontab{
		report:               report,
	}, nil

}

func (c *Crontab) Setup() error {
	crontab := cron.New()
	var err error

	// 每日0凌晨检查是否存在满足条件的归属于市场部的报单数据
	err = crontab.AddFunc("0 0 0 * * *", c.AutoChangeBelong)
	if err != nil {
		return err
	}

	crontab.Start()

	return nil
}

func (c *Crontab) Recover() {
	if p := recover(); p != nil {
		var buf [4096]byte
		n := runtime.Stack(buf[:], false)
		fmt.Printf("==> %s\n", string(buf[:n]))
		fmt.Println(p)
	}
}

