package server

import (
	"context"
	"github.com/XHXHXHX/medical_marketing/errs"
	"github.com/XHXHXHX/medical_marketing/service/customer_task"
	"github.com/XHXHXHX/medical_marketing/service/report"
	"github.com/XHXHXHX/medical_marketing/service/user"
	"github.com/XHXHXHX/medical_marketing/util/common"
	commonpb "github.com/XHXHXHX/medical_marketing_proto/gen/go/proto/common"
	"github.com/XHXHXHX/medical_marketing_proto/gen/go/proto/v1api"
	"time"
)

func (s *Server) ReportCreate(ctx context.Context, req *v1api.ReportCreateRequest) (*commonpb.Empty, error) {
	if req.GetConsumerMobile() == "" {
		return nil, errs.InvalidParams.Wrap("consumer_mobile")
	}
	if req.GetConsumerName() == "" {
		return nil, errs.InvalidParams.Wrap("consumer_name")
	}
	if req.GetExpectArriveTime() == 0 {
		return nil, errs.InvalidParams.Wrap("expire_arrive_time")
	}

	if common.PTimeUnix(req.GetExpectArriveTime()).Before(time.Now()) {
		return nil, errs.ExpectBeforeNow
	}

	err := s.reportService.Add(ctx, &report.Report{
		ReportUserID:     common.GetUserID(ctx),
		ConsumerMobile:   req.GetConsumerMobile(),
		ConsumerName:     req.GetConsumerName(),
		ExpectArriveTime: common.PTimeUnix(req.GetExpectArriveTime()),
		CreateTime:       common.PNow(),
	})
	if err != nil {
		return nil, err
	}
	return &commonpb.Empty{}, nil
}

func (s *Server) ReportRecover(ctx context.Context, req *v1api.ReportRecoverRequest) (*commonpb.Empty, error) {
	if req.GetId() == 0 {
		return nil, errs.InvalidParams.Wrap("id")
	}
	info, err := s.reportService.GetOne(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	if info.IsMatch.IsMatch() {
		return nil, errs.BanRecover
	}

	_, total, err := s.customerTask.List(ctx, &customer_task.SelectListRequest{
		ReportIDs:           []int64{info.ReportUserID},
	})
	if err != nil {
		return nil, err
	}
	if total > 0 {
		return nil, errs.BanRecover
	}

	return &commonpb.Empty{}, s.reportService.Del(ctx, req.GetId())
}

func (s *Server) ReportList(ctx context.Context, req *v1api.ReportListRequest) (*v1api.ReportListResponse, error) {
	page, err := common.GetPageInfo(req.GetPage())
	if err != nil {
		return nil, errs.InvalidParams.Wrap("page")
	}

	params := &report.SelectListRequest{
		UserId:    req.GetUserId(),
		CreateBeginTime: common.PTimeUnix(req.GetCreateStartTime()),
		CreateEndTime:   common.PTimeUnix(req.GetCreatEndTime()),
		ArriveStartTime: common.PTimeUnix(req.GetArriveStartTime()),
		ArriveEndTime: common.PTimeUnix(req.GetArriveEndTime()),
		IsMatch: report.IsMatch(req.IsMatch),
		Tag: req.GetTag(),
		Page:      page,
	}

	role := user.Role(common.GetRole(ctx))
	// 市场部的员工只能看自己的
	if role.IsMarketStaff() || role.IsCustomStaff() {
		params.UserId = common.GetUserID(ctx)
	}
	if role.IsCustomManager() {
		params.ShowCustomer = true
	}
	if role.IsMarket() {
		params.Belong = report.BelongMarket
	}
	if role.IsCustom() {
		params.Belong = report.BelongCustomer
	}

	if req.GetUserName() != "" {
		userList, _, err := s.userService.GetList(ctx, &user.SelectListRequest{
			Name:    req.GetUserName(),
			Status:  user.StatusNormal,
		})
		if err != nil {
			return nil, err
		}
		for _, v := range userList {
			params.UserIds = append(params.UserIds, v.ID)
		}
	}

	list, total, err := s.reportService.List(ctx, params)
	if err != nil {
		return nil, err
	}

	userIDs := make([]int64, 0, len(list))
	reportIDs := make([]int64, 0, len(list))
	// TODO 去重
	for _, v := range list {
		userIDs = append(userIDs, v.ReportUserID)
		reportIDs = append(reportIDs, v.ID)
	}

	var taskMap map[int64]*customer_task.CustomerTask

	if req.GetRelationTask() {
		taskList, _, err := s.customerTask.List(ctx, &customer_task.SelectListRequest{
			ReportIDs:           reportIDs,
		})
		if err != nil {
			return nil, err
		}

		taskMap = make(map[int64]*customer_task.CustomerTask)
		for _, v := range taskList {
			taskMap[v.ReportID] = v
			userIDs = append(userIDs, v.UserID)
		}
	}

	nameMap, err := s.userService.GetNameMap(ctx, userIDs)
	if err != nil {
		return nil, err
	}


	res := &v1api.ReportListResponse{
		List: make([]*commonpb.Report, 0, len(list)),
		Page: page.Page2Pagination(total),
	}

	for _, v := range list {
		tmp := &commonpb.Report{
			Id:                   v.ID,
			UserId:               v.ReportUserID,
			UserName:             nameMap[v.ReportUserID],
			ExceptArriveTime:     common.TimeToUnix(v.ExpectArriveTime),
			ConsumerMobile:       v.ConsumerMobile,
			ConsumerName:         v.ConsumerName,
			IsMatch:              int64(v.IsMatch),
			ActualArriveTime:     common.TimeToUnix(v.ActualArrivedTime),
			ConsumerAmount:       v.ConsumerAmount,
			CreateTime:           common.TimeToUnix(v.CreateTime),
			Tag:                  v.Tag,
			Memo:                 v.Memo,
		}
		if taskMap != nil && taskMap[v.ID] != nil {
			tmp.RelationTask = true
			tmp.RelationTaskUserId = taskMap[v.ID].UserID
			tmp.RelationTaskUsername = nameMap[taskMap[v.ID].UserID]
		}
		res.List = append(res.List, tmp)
	}

	return res, nil
}

func (s *Server) ReportChangeMatch(ctx context.Context, req *v1api.ReportChangeMatchRequest) (*commonpb.Empty, error) {
	if req.GetId() == 0 {
		return nil, errs.InvalidParams.Wrap("id")
	}

	match := report.IsMatch(req.GetIsMatch())
	if match.Valid() == false {
		return nil, errs.InvalidParams.Wrap("match")
	}

	info, err := s.reportService.GetOne(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	if info.IsMatch == match {
		return &commonpb.Empty{}, nil
	}

	info.IsMatch = match

	return &commonpb.Empty{}, s.reportService.Update(ctx, info)
}