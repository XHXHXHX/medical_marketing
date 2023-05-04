package report

import (
	"context"
)

type Service interface {
	List(ctx context.Context, req *SelectListRequest) ([]*Report, int64, error)
	Add(ctx context.Context, info *Report) error
	Del(ctx context.Context, id int64) error
	Import(ctx context.Context, buffer []byte) ([]*ImportErrorResult, error)
	SelectByConsumerMobile(ctx context.Context, mobile string) (*Report, error)
	GetOne(ctx context.Context, id int64) (*Report, error)
	GetMore(ctx context.Context, ids []int64) ([]*Report, error)
	GetMap(ctx context.Context, ids []int64) (map[int64]*Report, error)
	AutoChangeBelong(ctx context.Context)
	Update(ctx context.Context, info *Report) error
}
