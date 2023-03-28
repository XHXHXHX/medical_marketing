package server

import (
	"context"
	"github.com/XHXHXHX/medical_marketing/errs"
	"github.com/XHXHXHX/medical_marketing/service/customer_task"
	"github.com/XHXHXHX/medical_marketing/service/report"
	"github.com/XHXHXHX/medical_marketing/service/user"
	"github.com/XHXHXHX/medical_marketing/util/common"
	"github.com/XHXHXHX/medical_marketing_proto/gen/go/proto/v1api"
)

func (s *Server) StatisticsMarket(ctx context.Context, req *v1api.StatisticsMarketRequest) (*v1api.StatisticsMarketResponse, error) {
	page, err := common.GetPageInfo(req.GetPage())
	if err != nil {
		return nil, errs.InvalidParams.Wrap("page")
	}

	userList, total, err := s.userService.GetList(ctx, &user.SelectListRequest{
		Name: req.GetName(),
		Mobile: req.GetMobile(),
		Role: user.RoleMarketStaff,
		Page: page,
	})
	if err != nil {
		return nil, err
	}
	if len(userList) == 0 {
		return &v1api.StatisticsMarketResponse{
			Page: page.Page2Pagination(total),
		}, nil
	}

	userIds := make([]int64, 0, len(userList))
	mMap := make(map[int64]*v1api.StatisticsMarketResponse_Info, len(userList))
	for _, v := range userList {
		userIds = append(userIds, v.ID)
		mMap[v.ID] = &v1api.StatisticsMarketResponse_Info{
			UserName:   v.Name,
			UserMobile: v.Mobile,
			TotalNum:   0,
			FinishNum:  0,
		}
	}

	reportList, _, err := s.reportService.List(ctx, &report.SelectListRequest{
		UserIds:   userIds,
	})

	for _, v := range reportList {
		if mMap[v.ReportUserID] == nil {
			continue
		}

		mMap[v.ReportUserID].TotalNum++

		if v.IsMatch == report.Match {
			mMap[v.ReportUserID].FinishNum++
		}
	}

	res := &v1api.StatisticsMarketResponse{
		List: make([]*v1api.StatisticsMarketResponse_Info, 0, len(userList)),
		Page: page.Page2Pagination(total),
	}

	for _, v := range userList {
		res.List = append(res.List, mMap[v.ID])
	}

	return res, nil
}

func (s *Server) StatisticsCustomer(ctx context.Context, req *v1api.StatisticsCustomerRequest) (*v1api.StatisticsCustomerResponse, error) {
	page, err := common.GetPageInfo(req.GetPage())
	if err != nil {
		return nil, errs.InvalidParams.Wrap("page")
	}

	userList, total, err := s.userService.GetList(ctx, &user.SelectListRequest{
		Name: req.GetName(),
		Mobile: req.GetMobile(),
		Role: user.RoleCustomStaff,
		Page: page,
	})
	if err != nil {
		return nil, err
	}
	if len(userList) == 0 {
		return &v1api.StatisticsCustomerResponse{
			Page: page.Page2Pagination(total),
		}, nil
	}

	userIds := make([]int64, 0, len(userList))
	mMap := make(map[int64]*v1api.StatisticsCustomerResponse_Info, len(userList))
	for _, v := range userList {
		userIds = append(userIds, v.ID)
		mMap[v.ID] = &v1api.StatisticsCustomerResponse_Info{
			UserName:   v.Name,
			UserMobile: v.Mobile,
			TotalNum:   0,
			FinishNum:  0,
		}
	}

	taskList, _, err := s.customerTask.List(ctx, &customer_task.SelectListRequest{
		UserIDs:   userIds,
	})

	for _, v := range taskList {
		if mMap[v.UserID] == nil {
			continue
		}

		mMap[v.UserID].TotalNum++

		if v.IsFinished.IsFinished() {
			mMap[v.UserID].FinishNum++
		}
	}

	res := &v1api.StatisticsCustomerResponse{
		List: make([]*v1api.StatisticsCustomerResponse_Info, 0, len(userList)),
		Page: page.Page2Pagination(total),
	}

	for _, v := range userList {
		res.List = append(res.List, mMap[v.ID])
	}

	return res, nil
}
