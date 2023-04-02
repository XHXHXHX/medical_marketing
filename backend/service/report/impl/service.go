package impl

import (
	"context"
	"errors"
	"github.com/XHXHXHX/medical_marketing/errs"
	reportRepo "github.com/XHXHXHX/medical_marketing/repository/report"
	"github.com/XHXHXHX/medical_marketing/service/report"
	"github.com/XHXHXHX/medical_marketing/util/excel"
	"strconv"
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
	return s.repo.SelectList(ctx, req)
}

func (s *service) GetOne(ctx context.Context, id int64) (*report.Report, error) {
	return s.repo.SelectById(ctx, id)
}

func (s *service) GetMore(ctx context.Context, ids []int64) ([]*report.Report, error) {
	return s.repo.SelectByIds(ctx, ids)
}

func (s *service) GetMap(ctx context.Context, ids []int64) (map[int64]*report.Report, error) {
	list, err := s.GetMore(ctx, ids)
	if err != nil {
		return nil, err
	}

	mMap := make(map[int64]*report.Report, len(list))
	for _, v := range list {
		mMap[v.ID] = v
	}
	return mMap, nil
}

func (s *service) SelectByConsumerMobile(ctx context.Context, mobile string) (*report.Report, error) {
	return s.repo.SelectByMobile(ctx, mobile)
}

func (s *service) Add(ctx context.Context, info *report.Report) error {
	exist, err := s.repo.SelectByMobile(ctx, info.ConsumerMobile)
	if err != nil && errors.Is(err, errs.NotFoundData) == false {
		return err
	}

	if exist != nil {
		return errs.ExistSameConsumerMobile
	}

	info.IsMatch = report.UnMatch

	return s.repo.Insert(ctx, info)
}

func (s *service) Del(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) Import(ctx context.Context, buffer []byte) ([]*report.ImportErrorResult, error) {
	data, err := excel.ReadExcelForMap(buffer, map[string]int{
		"name":       	 0,
		"mobile":        1,
		"arrive_time":   2,
		"consume_amount": 3,
	}, 1)
	if err != nil {
		return nil, err
	}

	mobiles := make([]string, 0, len(data))
	for _, v := range data {
		mobiles = append(mobiles, v["mobile"])
	}

	reportList, _, err := s.repo.SelectList(ctx, &report.SelectListRequest{
		ConsumerMobiles: mobiles,
	})
	if err != nil {
		return nil, err
	}

	reportMap := make(map[string]*report.Report, len(reportList))
	for _, v := range reportList {
		reportMap[v.ConsumerMobile] = v
	}

	result := make([]*report.ImportErrorResult, 0)

	for i, v := range data {
		consumerName := v["name"]
		arriveTime := excel.ExcelDate(v["arrive_time"], "2006-01-02 15:04")
		consumeAmountF, err := strconv.ParseFloat(v["consume_amount"], 10)
		if err != nil {
			result = append(result, &report.ImportErrorResult{
				No: i+1,
				Error: "时间格式错误，eg: 2006-01-02 15:04:05",
			})
			continue
		}

		r, ok := reportMap[v["mobile"]]

		if !ok {
			result = append(result, &report.ImportErrorResult{
				No: i+1,
				Error: "未匹配",
			})
			continue
		}

		if r.IsMatch.IsMatch() {
			result = append(result, &report.ImportErrorResult{
				No: i+1,
				Error: "已匹配",
			})
			continue
		}

		r.ActualArrivedTime = arriveTime
		r.IsMatch = 1
		r.ConsumerAmount = int64(consumeAmountF * 100)
		r.ConsumerName = consumerName

		err = s.repo.Update(ctx, r)
		if err != nil {
			result = append(result, &report.ImportErrorResult{
				No: i+1,
				Error: "更新错误",
			})
			continue
		}
	}

	return result, nil
}
