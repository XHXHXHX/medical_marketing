package user

import (
	"time"

	"github.com/XHXHXHX/medical_marketing/util/common"
)

const (
	InitPassword = "A1234567"
	SaltLength = 4
	TokenLength = 128
	TokenExpire = 2 * time.Hour
	TokenPrefix = "token_%d"
)

type User struct {
	ID int64 `bson:"id"`
	Name string `bson:"name"`
	Role Role `bson:"role"`
	IsSupperAdmin int64 `bson:"is_supper_admin"`
	Mobile string `bson:"mobile"`
	Password string `bson:"password"`
	Salt string `bson:"salt"`
	Status Status `bson:"status"`
	CreateTime time.Time `bson:"create_time"`
}

func (u *User) IsAdmin() bool {
	return u.IsSupperAdmin == 1
}

const (
	RoleMarketManager Role = 1 //  市场部经理
	RoleMarketStaff Role = 2 //  市场部员工
	RoleCustomManager Role = 3 // 客服部经理
	RoleCustomStaff Role = 4 // 客服部员工
)

type Role int64

func (r Role) IsValid() bool {
	return r.IsMarket() || r.IsCustom()
}

func (r Role) IsMarketManager() bool {
	return r == RoleMarketManager
}

func (r Role) IsMarketStaff() bool {
	return r == RoleMarketStaff
}

func (r Role) IsCustomManager() bool {
	return r == RoleCustomManager
}

func (r Role) IsCustomStaff() bool {
	return r == RoleCustomStaff
}

func (r Role) IsMarket() bool {
	return r == RoleMarketManager || r == RoleMarketStaff
}

func (r Role) IsCustom() bool {
	return r == RoleCustomManager || r == RoleCustomStaff
}

func (r Role) IsManager() bool {
	return r == RoleMarketManager || r == RoleCustomManager
}

func (r Role) IsStaff() bool {
	return r == RoleMarketStaff || r == RoleCustomStaff
}

const (
	StatusInvalid Status = 0
	StatusNormal Status = 1
	StatusLeave Status = 2
)

type Status int64

func (s Status) IsValid() bool {
	return s.IsNormal() || s.IsLeave()
}
func (s Status) IsNormal() bool {
	return s == StatusNormal
}

func (s Status) IsLeave() bool {
	return s == StatusLeave
}

type SelectListRequest struct{
	Mobile string
	Mobiles []string
	Name string
	Role Role
	Roles []Role
	Status Status
	Page *common.Page
}