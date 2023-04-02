package customer_task

import (
	"context"
	"errors"
	"github.com/XHXHXHX/medical_marketing/errs"

	"github.com/XHXHXHX/medical_marketing/service/customer_task"

	"gorm.io/gorm"
)

type Repository interface {
	GetClient(ctx context.Context) *gorm.DB
	Begin(ctx context.Context) context.Context
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Insert(ctx context.Context, insert ...*customer_task.CustomerTask) error
	SelectById(ctx context.Context, id int64) (*customer_task.CustomerTask, error)
	SelectByIds(ctx context.Context, ids []int64) ([]*customer_task.CustomerTask, error)
	SelectList(ctx context.Context, req *customer_task.SelectListRequest) ([]*customer_task.CustomerTask, int64, error)
	Update(ctx context.Context, info *customer_task.CustomerTask) error
}

type repo struct {
	baseRepo
}

func NewRepo(client *gorm.DB) Repository {
	return &repo{
		baseRepo{client: client},
	}
}

func (repo *repo) Insert(ctx context.Context, insert ...*customer_task.CustomerTask) error {
	return repo.GetClient(ctx).Create(insert).Error
}

func (repo *repo) SelectById(ctx context.Context, id int64) (*customer_task.CustomerTask, error) {
	var info customer_task.CustomerTask
	err := repo.GetClient(ctx).Where("id = ?", id).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.NotFoundData
	}

	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (repo *repo) SelectByIds(ctx context.Context, ids []int64) ([]*customer_task.CustomerTask, error) {
	var list []*customer_task.CustomerTask
	err := repo.GetClient(ctx).Where("id in ?", ids).Scan(&list).Error

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (repo *repo) SelectList(ctx context.Context, req *customer_task.SelectListRequest) ([]*customer_task.CustomerTask, int64, error) {
	tx := repo.GetClient(ctx)

	if len(req.UserIDs) > 0 {
		tx = tx.Where("user_id in ?", req.UserIDs)
	}
	if len(req.ReportIDs) > 0 {
		tx = tx.Where("report_id in ?", req.ReportIDs)
	}
	if req.DistributeStartTime != nil {
		tx = tx.Where("distribute_time >= ?", req.DistributeStartTime)
	}
	if req.DistributeEndTime != nil {
		tx = tx.Where("distribute_time <= ?", req.DistributeEndTime)
	}
	if req.FinishStartTime != nil {
		tx = tx.Where("finish_time >= ?", req.FinishStartTime)
	}
	if req.FinishEndTime != nil {
		tx = tx.Where("finish_time <= ?", req.FinishEndTime)
	}
	if req.IsFinished.IsValid() {
		tx = tx.Where("is_finished = ?", req.IsFinished)
	}

	var total int64
	err := tx.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if req.Page != nil {
		tx = tx.Offset(int(req.Page.PageSize*(req.Page.CurrentPage-1))).Limit(int(req.Page.PageSize))
	}

	var list []*customer_task.CustomerTask
	err = tx.Order("id desc").Scan(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (repo *repo) Update(ctx context.Context, info *customer_task.CustomerTask) error {
	return repo.GetClient(ctx).Updates(info).Error
}


