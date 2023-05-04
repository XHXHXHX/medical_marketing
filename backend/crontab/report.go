package crontab

import (
	"context"
	"time"

	"github.com/XHXHXHX/medical_marketing/util/common"
	"github.com/XHXHXHX/medical_marketing/util/logx"
)

func (c *Crontab) AutoChangeBelong() {
	defer c.Recover()

	ctx, cancel := context.WithTimeout(common.SetGlobalID(context.Background()), 5 * time.Minute)
	defer cancel()

	logx.Infof(ctx, "Crontab AutoChangeBelong")

	c.report.AutoChangeBelong(ctx)
}
