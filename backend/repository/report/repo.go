package report

import (
	"context"
	"errors"
	"github.com/XHXHXHX/medical_marketing/errs"
	reportService "github.com/XHXHXHX/medical_marketing/service/report"

	"gorm.io/gorm"
)

type Repository interface {
	Insert(ctx context.Context, insert ...*reportService.Report) error
	SelectById(ctx context.Context, id int64) (*reportService.Report, error)
	SelectByIds(ctx context.Context, ids []int64) ([]*reportService.Report, error)
	SelectList(ctx context.Context, req *reportService.SelectListRequest) ([]*reportService.Report, int64, error)
	SelectByMobile(ctx context.Context, mobile string) (*reportService.Report, error)
	Delete(ctx context.Context, ids ...int64) error
	SelectUnMatchList(ctx context.Context) ([]*reportService.Report, error)
	Update(ctx context.Context, info *reportService.Report) error
	UpdateBelong(ctx context.Context, ids []int64, belong reportService.Belong) error
}

type repo struct {
	baseRepo
}

func NewRepo(client *gorm.DB) Repository {
	return &repo{
		baseRepo{client: client},
	}
}

func (repo *repo) Insert(ctx context.Context, insert ...*reportService.Report) error {
	return repo.GetClient(ctx).Create(insert).Error
}

func (repo *repo) SelectById(ctx context.Context, id int64) (*reportService.Report, error) {
	var info reportService.Report
	err := repo.GetClient(ctx).Where("id = ?", id).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.NotFoundData
	}

	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (repo *repo) SelectByIds(ctx context.Context, ids []int64) ([]*reportService.Report, error) {
	var list []*reportService.Report
	err := repo.GetClient(ctx).Where("id in ?", ids).Scan(&list).Error

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (repo *repo) SelectList(ctx context.Context, req *reportService.SelectListRequest) ([]*reportService.Report, int64, error) {
	tx := repo.GetClient(ctx)

	if req.UserId > 0 {
		tx = tx.Where("report_user_id = ?", req.UserId)
	}
	if len(req.UserIds) > 0 {
		tx = tx.Where("report_user_id in ?", req.UserIds)
	}
	if len(req.CustomerNames) > 0 {
		tx = tx.Where("consumer_name in ?", req.CustomerNames)
	}
 	if req.CreateBeginTime != nil {
		tx = tx.Where("create_time >= ?", req.CreateBeginTime)
	}
	if req.CreateEndTime != nil {
		tx = tx.Where("create_time <= ?", req.CreateEndTime)
	}
	if req.ArriveStartTime != nil {
		tx = tx.Where("actual_arrived_time >= ?", req.ArriveStartTime)
	}
	if req.ArriveEndTime != nil {
		tx = tx.Where("actual_arrived_time <= ?", req.ArriveEndTime)
	}
	if req.Tag != "" {
		tx = tx.Where("tag = ?", req.Tag)
	}
	if req.IsMatch.Valid() {
		tx = tx.Where("is_match = ?", req.IsMatch)
	}
	if len(req.ConsumerMobiles) > 0 {
		tx = tx.Where("consumer_mobile in ?", req.ConsumerMobiles)
	}
	if req.Belong.IsMarket() || req.Belong.IsCustomer() {
		tx = tx.Where("belong = ?", req.Belong)
	}
	var total int64
	err := tx.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if req.Page != nil {
		tx = tx.Offset(int(req.Page.PageSize*(req.Page.CurrentPage-1))).Limit(int(req.Page.PageSize))
	}

	var list []*reportService.Report
	err = tx.Order("id desc").Scan(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (repo *repo) SelectUnMatchList(ctx context.Context) ([]*reportService.Report, error) {
	var list []*reportService.Report
	err := repo.GetClient(ctx).Where("is_match = 2").Scan(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (repo *repo) SelectByMobile(ctx context.Context, mobile string) (*reportService.Report, error) {
	var info reportService.Report
	err := repo.GetClient(ctx).Where("consumer_mobile = ?", mobile).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.NotFoundData
	}

	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (repo *repo) Delete(ctx context.Context, ids ...int64) error {
	return repo.GetClient(ctx).Where("id in ?", ids).Delete(&reportService.Report{}).Error
}

func (repo *repo) Update(ctx context.Context, info *reportService.Report) error {
	return repo.GetClient(ctx).Updates(info).Error
}

func (repo *repo) UpdateBelong(ctx context.Context, ids []int64, belong reportService.Belong) error {
	return repo.GetClient(ctx).Where("id in ?", ids).UpdateColumn("belong", belong).Error
}


