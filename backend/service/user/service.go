package user

import "context"

type Service interface {
	Create(ctx context.Context, info *User, password string) error
	ChangePasswd(ctx context.Context, id int64, oldPasswd, newPasswd string) error
	ChangeStatus(ctx context.Context, id int64, status Status) error
	GetNameMap(ctx context.Context, ids []int64) (map[int64]string, error)
	GetOne(ctx context.Context, id int64) (*User, error)
	GetMore(ctx context.Context, ids []int64) ([]*User, error)
	GetList(ctx context.Context, req *SelectListRequest) ([]*User, int64, error)
	Auth(ctx context.Context, token string) (*User, error)
	Login(ctx context.Context, mobile, passwd string) (string, error)
	Logout(ctx context.Context) error
	CheckRole(ctx context.Context, targetRole Role) bool
}
