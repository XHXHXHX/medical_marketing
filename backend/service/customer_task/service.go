package customer_task

import "context"

type Service interface {
	Create(ctx context.Context, info *CustomerTask) error
	BatchCreate(ctx context.Context, userID int64, reportIDs []int64) error
	Finish(ctx context.Context, taskID int64, desc string, isFinished IsFinished) error
	List(ctx context.Context, req *SelectListRequest) ([]*CustomerTask, int64, error)
	HistoryList(ctx context.Context, taskIDs []int64) ([]*CustomerTaskHistory, error)
	Update(ctx context.Context, info *CustomerTask) error
	GetOne(ctx context.Context, taskID int64) (*CustomerTask, error)
}
