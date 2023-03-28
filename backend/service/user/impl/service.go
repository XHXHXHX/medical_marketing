package impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/XHXHXHX/medical_marketing/util/logx"
	"time"

	"github.com/XHXHXHX/medical_marketing/errs"
	userRepo "github.com/XHXHXHX/medical_marketing/repository/user"
	"github.com/XHXHXHX/medical_marketing/service/user"
	"github.com/XHXHXHX/medical_marketing/util/common"

	goredis "github.com/go-redis/redis/v8"
)

type service struct {
	repo userRepo.Repository
	redisCli *goredis.Client
}

func NewService(repo userRepo.Repository, redisCli *goredis.Client) user.Service {
	return &service{
		repo: repo,
		redisCli: redisCli,
	}
}

func (s *service) Create(ctx context.Context, info *user.User, password string) error {

	exists, err := s.repo.SelectByMobile(ctx, info.Mobile)
	if err != nil && errors.Is(err, errs.NotFoundData) == false {
		return err
	}
	if exists != nil {
		return errs.MobileRepeat
	}

	info.IsSupperAdmin = 2 // 后插入的数据固定不是超管
	info.Salt = common.RandString(user.SaltLength)
	info.Password = common.MD5Password(password, info.Salt)
	info.CreateTime = time.Now()
	info.Status = user.StatusNormal

	return s.repo.Insert(ctx, info)
}

func (s *service) ChangePasswd(ctx context.Context, id int64, oldPasswd, newPasswd string) error {

	info, err := s.repo.SelectById(ctx, id)
	if err != nil {
		return err
	}

	if common.MD5Password(oldPasswd, info.Salt) != info.Password {
		return errs.UserPasswordError
	}

	info.Password = common.MD5Password(newPasswd, info.Salt)

	return s.repo.Update(ctx, info)
}

func (s *service) ChangeStatus(ctx context.Context, id int64, status user.Status) error {
	info, err := s.repo.SelectById(ctx, id)
	if err != nil {
		return err
	}

	fmt.Println("check", s.CheckRole(ctx, info.Role))
	if s.CheckRole(ctx, info.Role) == false {
		return errs.RoleError
	}

	info.Status = status

	return s.repo.Update(ctx, info)
}

func (s *service) GetOne(ctx context.Context, id int64) (*user.User, error) {
	return s.repo.SelectById(ctx, id)
}

func (s *service) GetMore(ctx context.Context, ids []int64) ([]*user.User, error) {
	return s.repo.SelectByIds(ctx, ids)
}


func (s *service) GetNameMap(ctx context.Context, ids []int64) (map[int64]string, error) {
	list, err := s.repo.SelectByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	nMap := make(map[int64]string, len(list))
	for _, v := range list {
		nMap[v.ID] = v.Name
	}

	return nMap, err
}

func (s *service) GetList(ctx context.Context, req *user.SelectListRequest) ([]*user.User, int64, error) {
	return s.repo.SelectList(ctx, req)
}

func (s *service) Auth(ctx context.Context, token string) (*user.User, error) {

	id, err := s.redisCli.Get(ctx, token).Int64()
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}
	if id == 0 {
		return nil, errs.InvalidToken
	}

	info, err := s.repo.SelectById(ctx, id)
	if err != nil {
		return nil, errors.New("没有找到用户")
	}

	// 延长时效
	err = s.redisCli.SetEX(ctx, token, id, user.TokenExpire).Err()
	if err != nil {
		logx.Warnf(ctx, "延长token失效 [id:%d]", id)
	}

	// 延长时效
	err = s.redisCli.SetEX(ctx, s.getRedisTokenUserKey(id), token, user.TokenExpire).Err()
	if err != nil {
		logx.Warnf(ctx, "延长token失效 [id:%d]", id)
	}

	return info, nil
}

func (s *service) Login(ctx context.Context, mobile, passwd string) (string, error) {
	info, err := s.repo.SelectByMobile(ctx, mobile)
	if err != nil {
		return "", err
	}

	if common.MD5Password(passwd, info.Salt) != info.Password {
		return "", errs.UserPasswordError
	}

	ctx = common.SetUserID(ctx, info.ID)
	ctx = common.SetUserName(ctx, info.Name)
	ctx = common.SetRole(ctx, int64(info.Role))

	return s.setToken(ctx, info.ID)
}

func (s *service) Logout(ctx context.Context) error {
	token := s.getTokenFromUserID(ctx, common.GetUserID(ctx))

	if token == "" {
		return nil
	}

	defer s.redisCli.Del(ctx, s.getRedisTokenUserKey(common.GetUserID(ctx)))

	return s.redisCli.Del(ctx, token).Err()
}

func (s *service) getTokenFromUserID(ctx context.Context, id int64) string {
	return s.redisCli.Get(ctx, s.getRedisTokenUserKey(id)).Val()
}

func (s *service) setToken(ctx context.Context, id int64) (string, error) {
	token := common.RandNormalString(user.TokenLength)

	oldToken := s.getTokenFromUserID(ctx, id)

	err := s.redisCli.SetEX(ctx, token, id, user.TokenExpire).Err()
	if err != nil {
		return "", err
	}

	if oldToken != "" {
		s.redisCli.Del(ctx, oldToken)
	}

	err = s.redisCli.SetEX(ctx, s.getRedisTokenUserKey(id), token, user.TokenExpire).Err()
	if err != nil {
		s.redisCli.Del(ctx, token)
		return "", err
	}

	return token, nil
}

func (s *service) getRedisTokenUserKey(id int64) string {
	return fmt.Sprintf(user.TokenPrefix, id)
}

func (s *service) CheckRole(ctx context.Context, targetRole user.Role) bool {
	if common.GetAdmin(ctx) {
		return true
	}

	role := user.Role(common.GetRole(ctx))
	if role.IsStaff() {
		return false
	}

	if role.IsMarket() && targetRole.IsMarket() {
		return true
	}

	if role.IsCustom() && targetRole.IsCustom() {
		return true
	}

	return false
}