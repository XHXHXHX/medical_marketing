package report

import (
	"context"

	"github.com/XHXHXHX/medical_marketing/repository"

	"gorm.io/gorm"
)

const (
	TableName = "user"
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
