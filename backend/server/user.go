package server

import (
	"context"
	"github.com/XHXHXHX/medical_marketing/errs"
	"github.com/XHXHXHX/medical_marketing/service/user"
	"github.com/XHXHXHX/medical_marketing/util/common"
	commonpb "github.com/XHXHXHX/medical_marketing_proto/gen/go/proto/common"
	"github.com/XHXHXHX/medical_marketing_proto/gen/go/proto/v1api"
)


func (s *Server) Login(ctx context.Context, req *v1api.LoginRequest) (*v1api.LoginResponse, error) {
	if req.GetMobile() == "" {
		return nil, errs.InvalidParams.Wrap("mobile")
	}
	if req.GetPassword() == "" {
		return nil, errs.InvalidParams.Wrap("password")
	}

	token, err := s.userService.Login(ctx, req.GetMobile(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &v1api.LoginResponse{
		Token: token,
	}, nil
}

func (s *Server) Logout(ctx context.Context, req *commonpb.Empty) (*commonpb.Empty, error) {
	err := s.userService.Logout(ctx)
	if err != nil {
		return nil, err
	}

	return &commonpb.Empty{}, nil
}

func (s *Server) UserInfo(ctx context.Context, req *commonpb.Empty) (*v1api.UserInfoResponse, error) {
	userInfo, err := s.userService.GetOne(ctx, common.GetUserID(ctx))
	if err != nil {
		return nil, err
	}
	return &v1api.UserInfoResponse{
		UserId:     userInfo.ID,
		Name:   userInfo.Name,
		Mobile: userInfo.Mobile,
		Status:     int64(userInfo.Status),
		IsAdmin: userInfo.IsAdmin(),
		Role: int64(userInfo.Role),
	}, nil
}

func (s *Server) UserChangePasswd(ctx context.Context, req *v1api.UserChangePasswdRequest) (*commonpb.Empty, error) {
	if req.GetOldPasswd() == "" {
		return nil, errs.InvalidParams.Wrap("old_password")
	}
	if req.GetNewPasswd() == "" {
		return nil, errs.InvalidParams.Wrap("new_password")
	}
	err := s.userService.ChangePasswd(ctx, common.GetUserID(ctx), req.GetOldPasswd(), req.GetNewPasswd())
	if err != nil {
		return nil, err
	}

	return &commonpb.Empty{}, err
}

func (s *Server) UserCreate(ctx context.Context, req *v1api.UserCreateRequest) (*v1api.UserCreateResponse, error) {
	if req.GetMobile() == "" {
		return nil, errs.InvalidParams.Wrap("mobile")
	}
	if req.GetName() == "" {
		return nil, errs.InvalidParams.Wrap("name")
	}
	if user.Role(req.GetRole()).IsValid() == false {
		return nil, errs.InvalidParams.Wrap("role")
	}

	err := s.userService.Create(ctx, &user.User{
		Name:          req.GetName(),
		Role:          user.Role(req.GetRole()),
		Mobile:        req.GetMobile(),
	}, user.InitPassword)
	if err != nil {
		return nil, err
	}

	return &v1api.UserCreateResponse{Password: user.InitPassword}, nil
}

// 暂时不要
func (s *Server) UserChangeRole(ctx context.Context, req *v1api.UserChangeRoleRequest) (*commonpb.Empty, error) {
	if req.GetUserId() == 0 {
		return nil, errs.InvalidParams.Wrap("user_id")
	}
	if user.Role(req.GetRole()).IsValid() ==  false {
		return nil, errs.InvalidParams.Wrap("role")
	}

	if s.userService.CheckRole(ctx, user.Role(req.GetRole())) == false {
		return nil, errs.InvalidParams
	}

	return &commonpb.Empty{}, nil
}

// 离职
func (s *Server) UserChangeStatus(ctx context.Context, req *v1api.UserChangeStatusRequest) (*commonpb.Empty, error) {
	if req.GetUserId() == 0 {
		return nil, errs.InvalidParams.Wrap("user_id")
	}

	err := s.userService.ChangeStatus(ctx, req.GetUserId(), user.StatusLeave)
	if err != nil {
		return nil, err
	}

	return &commonpb.Empty{}, nil
}

func (s *Server) UserList(ctx context.Context, req *v1api.UserListRequest) (*v1api.UserListResponse, error) {
	page, err := common.GetPageInfo(req.GetPage())
	if err != nil {
		return nil, errs.InvalidParams.Wrap("page")
	}

	list, total, err := s.userService.GetList(ctx, &user.SelectListRequest{
		Mobile:  req.GetMobile(),
		Mobiles: nil,
		Name:    req.GetName(),
		Role:    user.Role(req.GetRole()),
		Roles:   nil,
		Status:  user.StatusNormal,
		Page:    page,
	})
	if err != nil {
		return nil, err
	}

	res := &v1api.UserListResponse{
		Page: page.Page2Pagination(total),
		List: make([]*v1api.UserListResponse_User, 0, len(list)),
	}

	for _, v := range list {
		res.List = append(res.List, &v1api.UserListResponse_User{
			UserId: v.ID,
			Name:   v.Name,
			Mobile: v.Mobile,
			Role:   int64(v.Role),
			Status: int64(v.Status),
		})
	}

	return res, nil
}