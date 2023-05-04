package impl

import (
	"context"
	"errors"
	"github.com/XHXHXHX/medical_marketing/errs"
	reportRepo "github.com/XHXHXHX/medical_marketing/repository/report"
	"github.com/XHXHXHX/medical_marketing/service/report"
	"github.com/XHXHXHX/medical_marketing/util/common"
	"github.com/XHXHXHX/medical_marketing/util/excel"
	"github.com/XHXHXHX/medical_marketing/util/logx"
	"strconv"
	"time"
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

	names := make([]string, 0, len(data))
	for _, v := range data {
		names = append(names, v["name"])
	}

	reportList, _, err := s.repo.SelectList(ctx, &report.SelectListRequest{
		CustomerNames: names,
	})
	if err != nil {
		return nil, err
	}

	reportMap := make(map[string]*report.Report, len(reportList))
	for _, v := range reportList {
		reportMap[v.ConsumerName] = v
	}

	result := make([]*report.ImportErrorResult, 0)

	for i, v := range data {
		consumerMobile := v["mobile"]
		consumerName := v["name"]
		arriveTime := excel.ExcelDate(v["arrive_time"], "2006-01-02 15:04")
		consumeAmountF, err := strconv.ParseFloat(v["consume_amount"], 10)
		if err != nil {
			result = append(result, &report.ImportErrorResult{
				No: i+1,
				Error: "金额错误",
			})
			continue
		}
		consumeAmount := int64(consumeAmountF * 100)

		r, ok := reportMap[consumerName]

		if !ok {
			result = append(result, &report.ImportErrorResult{
				No: i+1,
				Error: "姓名未匹配",
			})
			continue
		}

		if r.ConsumerMobile == consumerName && r.ConsumerAmount == consumeAmount &&  r.IsMatch.IsMatch() {
			result = append(result, &report.ImportErrorResult{
				No: i+1,
				Error: "已匹配",
			})
			continue
		}

		r.ActualArrivedTime = arriveTime
		r.IsMatch = 1
		r.ConsumerAmount = consumeAmount
		r.ConsumerMobile = consumerMobile

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

func (s *service) Update(ctx context.Context, info *report.Report) error {
	return s.repo.Update(ctx, info)
}

/*
 报单7日后未自动匹配的数据，归属于客服
 到访三个月后的数据，归属于客服
 */
func (s *service) AutoChangeBelong(ctx context.Context) {

	// 1. 报单7日后未自动匹配的数据，归属于客服
	s.unMatchDataDistributeCustomer(ctx)
	// 2. 到访三个月后的数据，归属于客服
	s.matchedDataDistributeCustomer(ctx)
}

// 1. 报单7日后未自动匹配的数据，归属于客服
func (s *service) unMatchDataDistributeCustomer(ctx context.Context) {
	req := &report.SelectListRequest{
		IsMatch: report.UnMatch,
		Belong: report.BelongMarket,
		CreateEndTime: common.PTime(time.Now().Add(report.UnMatchDataDistributeCustomerDate)),
	}

	list, _, err := s.repo.SelectList(ctx, req)
	if err != nil {
		logx.Errorf(ctx, "unMatchDataDistributeCustomer select error", err)
	}

	err = s.changeBelongToCustomer(ctx, list)
	if err != nil {
		logx.Errorf(ctx, "unMatchDataDistributeCustomer changeBelongToCustomer error", err)
	}
}

// 2. 到访三个月后的数据，归属于客服
func (s *service) matchedDataDistributeCustomer(ctx context.Context) {
	req := &report.SelectListRequest{
		IsMatch: report.Match,
		Belong: report.BelongMarket,
		CreateEndTime: common.PTime(time.Now().Add(report.MatchedDataDistributeCustomerDate)),
	}

	list, _, err := s.repo.SelectList(ctx, req)
	if err != nil {
		logx.Errorf(ctx, "matchedDataDistributeCustomer select error", err)
	}

	err = s.changeBelongToCustomer(ctx, list)
	if err != nil {
		logx.Errorf(ctx, "matchedDataDistributeCustomer changeBelongToCustomer error", err)
	}
}

func (s *service) changeBelongToCustomer(ctx context.Context, list []*report.Report) error {
	ids := make([]int64, 0, len(list))
	for _, v := range list {
		ids = append(ids, v.ID)
	}

	return s.repo.UpdateBelong(ctx, ids, report.BelongCustomer)
}