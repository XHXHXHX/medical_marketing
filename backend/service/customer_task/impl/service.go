package impl

import (
	"context"
	"errors"
	customerTaskRepo "github.com/XHXHXHX/medical_marketing/repository/customer_task"
	customerTaskHistoryRepo "github.com/XHXHXHX/medical_marketing/repository/customer_task_history"
	"github.com/XHXHXHX/medical_marketing/service/customer_task"
	"github.com/XHXHXHX/medical_marketing/util/common"
	"time"
)

type service struct {
	customerTaskRepo customerTaskRepo.Repository
	customerTaskHistoryRepo customerTaskHistoryRepo.Repository
}

func NewService(customerTaskRepo customerTaskRepo.Repository, customerTaskHistoryRepo customerTaskHistoryRepo.Repository) customer_task.Service {
	return &service{
		customerTaskRepo: customerTaskRepo,
		customerTaskHistoryRepo: customerTaskHistoryRepo,
	}
}

func (s *service) Create(ctx context.Context, info *customer_task.CustomerTask) error {
	info.DistributeTime = time.Now()
	info.IsFinished = customer_task.UnFinish
	return s.customerTaskRepo.Insert(ctx, info)
}

func (s *service) GetOne(ctx context.Context, taskID int64) (*customer_task.CustomerTask, error) {
	return s.customerTaskRepo.SelectById(ctx, taskID)
}

func (s *service) BatchCreate(ctx context.Context, userID int64, reportIDs []int64) error {
	t := time.Now()
	insert := make([]*customer_task.CustomerTask, 0, len(reportIDs))
	for _, reportID := range reportIDs {
		insert = append(insert, &customer_task.CustomerTask{
			ReportID:       reportID,
			UserID:         userID,
			DistributeTime: t,
		})
	}

	return s.customerTaskRepo.Insert(ctx, insert...)
}

func (s *service) Finish(ctx context.Context, taskID int64, desc string, isFinished customer_task.IsFinished) error {
	info, err := s.customerTaskRepo.SelectById(ctx, taskID)
	if err != nil {
		return err
	}

	if info.IsFinished.IsFinished() {
		return errors.New("任务已经完成")
	}

	t := time.Now()

	err = s.customerTaskHistoryRepo.Insert(ctx, &customer_task.CustomerTaskHistory{
		UserID: common.GetUserID(ctx),
		TaskID:     taskID,
		Desc:       desc,
		CreateTime: t,
	})
	if err != nil {
		return err
	}

	if isFinished.IsFinished() == false {
		return nil
	}

	info.IsFinished = isFinished
	info.FinishTime = &t

	return s.customerTaskRepo.Update(ctx, info)
}

func (s *service) Update(ctx context.Context, info *customer_task.CustomerTask) error {
	return s.customerTaskRepo.Update(ctx, info)
}

func (s *service) List(ctx context.Context, req *customer_task.SelectListRequest) ([]*customer_task.CustomerTask, int64, error) {
	return s.customerTaskRepo.SelectList(ctx, req)
}

func (s *service) HistoryList(ctx context.Context, taskIDs []int64) ([]*customer_task.CustomerTaskHistory, error) {
	return s.customerTaskHistoryRepo.SelectByReportIds(ctx, taskIDs...)
}