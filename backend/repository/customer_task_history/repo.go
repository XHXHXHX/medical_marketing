package customer_task_history

import (
	"context"
	"errors"
	"github.com/XHXHXHX/medical_marketing/errs"

	"github.com/XHXHXHX/medical_marketing/service/customer_task"

	"gorm.io/gorm"
)

type Repository interface {
	Insert(ctx context.Context, insert ...*customer_task.CustomerTaskHistory) error
	SelectByReportIds(ctx context.Context, taskIDs ...int64) ([]*customer_task.CustomerTaskHistory, error)
	SelectLastOne(ctx context.Context, taskID int64) (*customer_task.CustomerTaskHistory, error)
}

type repo struct {
	baseRepo
}

func NewRepo(client *gorm.DB) Repository {
	return &repo{
		baseRepo{client: client},
	}
}

func (repo *repo) Insert(ctx context.Context, insert ...*customer_task.CustomerTaskHistory) error {
	return repo.GetClient(ctx).Create(insert).Error
}

func (repo *repo) SelectByReportIds(ctx context.Context, taskIDs ...int64) ([]*customer_task.CustomerTaskHistory, error) {
	var list []*customer_task.CustomerTaskHistory
	err := repo.GetClient(ctx).Where("task_id in ?", taskIDs).Order("id desc").Scan(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (repo *repo) SelectLastOne(ctx context.Context, taskID int64) (*customer_task.CustomerTaskHistory, error) {
	var info customer_task.CustomerTaskHistory
	err := repo.GetClient(ctx).Where("task_id = ?", taskID).Last(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.NotFoundData
	}

	if err != nil {
		return nil, err
	}

	return &info, nil
}


