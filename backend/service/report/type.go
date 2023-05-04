package report

import (
	"github.com/XHXHXHX/medical_marketing/util/common"
	"time"
)

const (
	OneDay = 24 * time.Hour
	UnMatchDataDistributeCustomerDate = 7 * OneDay // 报单7日后未自动匹配的数据，归属于客服
	MatchedDataDistributeCustomerDate = 3 * 30 * OneDay // 到访三个月后的数据，归属于客服
)

const (
	Match IsMatch = 1
	UnMatch IsMatch = 2
)

type IsMatch int64

func(m IsMatch) Valid() bool {
	return m == Match || m == UnMatch
}

func (m IsMatch) IsMatch() bool {
	return m == Match
}

// 归属
type Belong int64

const (
	BelongAll = 0
	BelongMarket = 1
	BelongCustomer = 2
)

func (b Belong) IsAll() bool {
	return b == BelongAll
}

func (b Belong) IsMarket() bool {
	return b == BelongMarket
}

func (b Belong) IsCustomer() bool {
	return b == BelongCustomer
}

type SelectListRequest struct {
	CustomerNames []string
	UserId int64
	UserIds []int64
	IsMatch IsMatch
	ConsumerMobiles []string
	ShowCustomer bool
	Belong Belong
	CreateBeginTime *time.Time
	CreateEndTime *time.Time
	ArriveStartTime *time.Time
	ArriveEndTime *time.Time
	Tag string
	Page *common.Page
}

type Report struct {
	ID int64 `bson:"id"`
	ReportUserID int64 `bson:"report_user_id"` // 报告人ID
	ConsumerMobile string `bson:"consumer_mobile"` // 客户手机号
	ConsumerName string `bson:"consumer_name"` // 客户姓名
	ConsumerAmount int64 `bson:"consumer_amount"` // 消费金额
	ExpectArriveTime *time.Time `bson:"expect_arrive_time"` // 预期到访日期
	IsMatch IsMatch `bson:"is_match"` // 是否到访 1是2否
	ActualArrivedTime *time.Time `bson:"actual_arrived_time"` // 到访时间
	PatientID int64 `bson:"patient_id"` // 患者ID
	CreateTime *time.Time `bson:"create_time"`
	IsDeleted int32 `bson:"is_deleted"` // 是否删除 1是 0 否
	DeleteTime *time.Time `bson:"delete_time"`
	Belong Belong `bson:"belong"` // 归属
	Tag string `bson:"tag"`
	Memo string `bson:"memo"`
}

type ImportErrorResult struct {
	No int
	Error string
}