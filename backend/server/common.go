package server

import (
	"github.com/XHXHXHX/medical_marketing/service/report"
)

type Common struct {
	reportService report.Service
}

func NewCommon(reportService report.Service) *Common {
	return &Common{
		reportService: reportService,
	}
}
