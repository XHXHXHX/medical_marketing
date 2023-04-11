package report

import (
	"context"
	"errors"
	"fmt"
	"github.com/XHXHXHX/medical_marketing/errs"

	userService "github.com/XHXHXHX/medical_marketing/service/user"

	"gorm.io/gorm"
)

type Repository interface {
	Insert(ctx context.Context, insert ...*userService.User) error
	SelectById(ctx context.Context, id int64) (*userService.User, error)
	SelectByIds(ctx context.Context, ids []int64) ([]*userService.User, error)
	SelectList(ctx context.Context, req *userService.SelectListRequest) ([]*userService.User, int64, error)
	SelectByMobile(ctx context.Context, mobile string) (*userService.User, error)
	Update(ctx context.Context, info *userService.User) error
}

type repo struct {
	baseRepo
}

func NewRepo(client *gorm.DB) Repository {
	return &repo{
		baseRepo{client: client},
	}
}

func (repo *repo) Insert(ctx context.Context, insert ...*userService.User) error {
	return repo.GetClient(ctx).Create(insert).Error
}

func (repo *repo) SelectById(ctx context.Context, id int64) (*userService.User, error) {
	var info userService.User
	err := repo.GetClient(ctx).Where("id = ?", id).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.NotFoundData
	}

	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (repo *repo) SelectByIds(ctx context.Context, ids []int64) ([]*userService.User, error) {
	var list []*userService.User
	err := repo.GetClient(ctx).Where("id in ?", ids).Scan(&list).Error

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (repo *repo) SelectList(ctx context.Context, req *userService.SelectListRequest) ([]*userService.User, int64, error) {
	tx := repo.GetClient(ctx)

	if req.Name != "" {
		tx = tx.Where("name like ?", fmt.Sprintf("%%%s%%", req.Name))
	}

	if req.Mobile != "" {
		tx = tx.Where("mobile = ?", req.Mobile)
	}

	if len(req.Mobiles) > 0 {
		tx = tx.Where("mobile in ?", req.Mobiles)
	}

	if req.Role.IsValid() {
		tx = tx.Where("role = ?", req.Role)
	}
	if len(req.Roles) > 0 {
		tx = tx.Where("role in ?", req.Roles)
	}

	if req.Status.IsValid() {
		tx = tx.Where("status = ?", req.Status)
	}

	var total int64
	err := tx.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if req.Page != nil {
		tx = tx.Offset(int(req.Page.PageSize*(req.Page.CurrentPage-1))).Limit(int(req.Page.PageSize))
	}

	var list []*userService.User
	err = tx.Order("status asc, id desc").Scan(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (repo *repo) SelectByMobile(ctx context.Context, mobile string) (*userService.User, error) {
	var info userService.User
	err := repo.GetClient(ctx).Where("mobile = ?", mobile).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.NotFoundData
	}

	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (repo *repo) Update(ctx context.Context, info *userService.User) error {
	return repo.GetClient(ctx).Updates(info).Error
}


