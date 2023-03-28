package repository

import (
	"context"
	"gorm.io/gorm"
)

type transaction struct{}

func GetClient(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(transaction{}).(*gorm.DB); ok {
		return tx
	}
	return nil
}

func SetClient(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, transaction{}, tx)
}
