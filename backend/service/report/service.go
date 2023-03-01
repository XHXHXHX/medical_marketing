package report

import (
	"context"
)

type Service interface {
	List(ctx context.Context, req *SelectListRequest) ([]*Report, int64, error)
	Add(ctx context.Context, info *Report) error
	Del(ctx context.Context, id int64) error
}
