package customer_task

import (
	"context"

	"github.com/XHXHXHX/medical_marketing/repository"

	"gorm.io/gorm"
)

const (
	TableName = "customer_task"
)


type baseRepo struct {
	client *gorm.DB
}

func (base *baseRepo) GetTableName() string {
	return TableName
}

func (base *baseRepo) GetClient(ctx context.Context) *gorm.DB {
	tx := repository.GetClient(ctx)
	if tx == nil {
		tx = base.client
	}
	return base.client.WithContext(ctx).Table(base.GetTableName())
}

func (base *baseRepo) Begin(ctx context.Context) context.Context {
	tx := base.GetClient(ctx)
	repository.SetClient(ctx, tx)
	return ctx
}

func (base *baseRepo) Commit(ctx context.Context) error {
	tx := base.GetClient(ctx)
	return tx.Commit().Error
}

func (base *baseRepo) Rollback(ctx context.Context) error {
	tx := base.GetClient(ctx)
	return tx.Rollback().Error
}