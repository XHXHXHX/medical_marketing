package impl

import (
	"context"
	"github.com/XHXHXHX/medical_marketing/errs"
	reportRepo "github.com/XHXHXHX/medical_marketing/repository/report"
	"github.com/XHXHXHX/medical_marketing/service/report"
)

type service struct {
	repo reportRepo.Repository
}

func NewService(repo reportRepo.Repository) report.Service {
	return &service{
		repo: repo,
	}
}

func (s *service) List(ctx context.Context, req *report.SelectListRequest) ([]*report.Report, int64, error) {
	// TODO 员工姓名查询
	return s.repo.SelectList(ctx, req)
}

func (s *service) Add(ctx context.Context, info *report.Report) error {
	exist, err := s.repo.SelectByMobile(ctx, info.ConsumerMobile)
	if err != nil {
		return err
	}

	if exist != nil {
		return errs.ExistSameConsumerMobile
	}

	return s.repo.Insert(ctx, info)
}

func (s *service) Del(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
