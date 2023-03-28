package server

import (
	"context"
	"errors"
	"github.com/XHXHXHX/medical_marketing/errs"
	"github.com/XHXHXHX/medical_marketing/service/customer_task"
	"github.com/XHXHXHX/medical_marketing/service/user"
	"github.com/XHXHXHX/medical_marketing/util/common"
	commonpb "github.com/XHXHXHX/medical_marketing_proto/gen/go/proto/common"
	"github.com/XHXHXHX/medical_marketing_proto/gen/go/proto/v1api"
	"time"
)

func(s *Server) CustomerServerDistribute(ctx context.Context, req *v1api.CustomerDistributeRequest) (*commonpb.Empty, error) {
	if req.GetUserId() == 0 {
		return nil, errs.InvalidParams.Wrap("user_id")
	}
	if len(req.GetReportIds()) == 0 {
		return nil, errs.InvalidToken.Wrap("report_ids")
	}

	role := user.Role(common.GetRole(ctx))

	if common.GetAdmin(ctx) == false && role.IsCustomManager() == false {
		return nil, errs.NoRole
	}

	userInfo, err := s.userService.GetOne(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	if userInfo.Role.IsCustomStaff() == false {
		return nil, errors.New("只能给客服分配任务")
	}

	taskList, _, err := s.customerTask.List(ctx, &customer_task.SelectListRequest{
		ReportIDs:           req.GetReportIds(),
	})
	if err != nil {
		return nil, err
	}

	existsTask := make(map[int64]int, len(taskList))
	for _, v := range taskList {
		existsTask[v.ReportID]++
		if v.UserID == req.GetUserId() {
			continue
		}
		if v.IsFinished.IsFinished() { // 已经完成的也不需要分配
			continue
		}
		v.UserID = req.GetUserId()
		v.DistributeTime = time.Now()
		err = s.customerTask.Update(ctx, v)
		if err != nil {
			return nil, err
		}
	}

	insertReport := make([]int64, 0, len(req.GetReportIds()))
	for _, v := range req.GetReportIds() {
		if existsTask[v] > 0 {
			continue
		}
		insertReport = append(insertReport, v)
	}

	if len(insertReport) == 0 {
		return &commonpb.Empty{}, nil
	}

	return &commonpb.Empty{}, s.customerTask.BatchCreate(ctx, req.GetUserId(), insertReport)
}

func (s *Server) CustomerServerList(ctx context.Context, req *v1api.CustomerReportListRequest) (*v1api.CustomerReportListResponse, error) {
	page, err := common.GetPageInfo(req.GetPage())
	if err != nil {
		return nil, errs.InvalidParams.Wrap("page")
	}

	params := &customer_task.SelectListRequest{
		IsFinished:          customer_task.IsFinished(req.GetIsFinished()),
		DistributeStartTime: common.PTimeUnix(req.GetDistributeStartTime()),
		DistributeEndTime:   common.PTimeUnix(req.GetDistributeEndTime()),
		FinishStartTime:     common.PTimeUnix(req.GetFinishStartTime()),
		FinishEndTime:       common.PTimeUnix(req.GetFinishEndTime()),
		Page:                page,
	}

	if req.GetMobile() != "" || req.GetName() != "" {
		list, _, err := s.userService.GetList(ctx, &user.SelectListRequest{
			Mobile:  req.GetMobile(),
			Name:    req.GetName(),
			Role:    user.RoleCustomStaff,
			Status:  user.StatusNormal,
		})
		if err != nil {
			return nil, err
		}

		for _, v := range list {
			params.UserIDs = append(params.UserIDs, v.ID)
		}
	}

	if req.GetCustomerMobile() != "" {
		rInfo, err := s.reportService.SelectByConsumerMobile(ctx, req.GetCustomerMobile())
		if err != nil {
			return nil, err
		}
		params.ReportIDs = append(params.ReportIDs, rInfo.ID)
	}

	list, total, err := s.customerTask.List(ctx, params)
	if err != nil {
		return nil, err
	}

	unique := make(map[int64]int, len(list))
	userIds := make([]int64, 0, len(list))
	reportIds := make([]int64, 0, len(list))
	taskIds := make([]int64, 0, len(list))

	for _, v := range list {
		unique[v.UserID]++
		userIds = append(userIds, v.UserID)
		reportIds = append(reportIds, v.ReportID)
		taskIds = append(taskIds, v.ID)
	}

	reportMap, err := s.reportService.GetMap(ctx, reportIds)
	if err != nil {
		return nil, err
	}

	history, err := s.customerTask.HistoryList(ctx, taskIds)
	if err != nil {
		return nil, err
	}

	for _, v := range history {
		if unique[v.UserID] > 0 {
			continue
		}
		unique[v.UserID]++
		userIds = append(userIds, v.UserID)
	}

	userNameMap, err := s.userService.GetNameMap(ctx, userIds)
	if err != nil {
		return nil, err
	}

	historyMap := make(map[int64][]*commonpb.CustomerServerHistory, len(list))
	for _, v := range history {
		historyMap[v.TaskID] = append(historyMap[v.TaskID], &commonpb.CustomerServerHistory{
			UserId: v.UserID,
			Name: userNameMap[v.UserID],
			TaskId: v.TaskID,
			Desc:       v.Desc,
			Time: common.TimeToUnix(&v.CreateTime),
		})
	}


	res := &v1api.CustomerReportListResponse{
		List: make([]*commonpb.CustomerServerInfo, 0, len(list)),
		Page: page.Page2Pagination(total),
	}

	for _, v := range list {
		lastDesc := ""
		if len(historyMap[v.ID]) > 0 {
			lastDesc = historyMap[v.ID][0].GetDesc()
		}
		res.List = append(res.List, &commonpb.CustomerServerInfo{
			TaskId:         v.ID,
			ReportId:       v.ReportID,
			UserId:         v.UserID,
			Name:           userNameMap[v.UserID],
			DistributeTime: common.TimeToUnix(&v.DistributeTime),
			CustomerName:   reportMap[v.ReportID].ConsumerName,
			CustomerAmount: reportMap[v.ReportID].ConsumerAmount,
			CustomerMobile: reportMap[v.ReportID].ConsumerMobile,
			IsFinished:     int64(v.IsFinished),
			History:        historyMap[v.ID],
			LastDesc:       lastDesc,
			FinishTime:     common.TimeToUnix(v.FinishTime),
		})
	}

	return res, nil
}

func (s *Server) CustomerServerResult(ctx context.Context, req *v1api.CustomerServerResultRequest) (*commonpb.Empty, error) {
	if req.GetDesc() == "" {
		return nil, errs.InvalidParams.Wrap("desc")
	}
	if req.GetTaskId() == 0 {
		return nil, errs.InvalidParams.Wrap("task_id")
	}
	if customer_task.IsFinished(req.GetIsFinished()).IsValid() == false {
		return nil, errs.InvalidParams.Wrap("is_finished")
	}

	task, err := s.customerTask.GetOne(ctx, req.GetTaskId())
	if err != nil {
		return nil, err
	}
	if task.UserID != common.GetUserID(ctx) {
		return nil, errs.NotBelongWithYou
	}

	err = s.customerTask.Finish(ctx, req.GetTaskId(), req.GetDesc(), customer_task.IsFinished(req.GetIsFinished()))
	if err != nil {
		return nil, err
	}
	return &commonpb.Empty{}, nil
}