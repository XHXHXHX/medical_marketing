package report

import (
	"context"
	"github.com/XHXHXHX/medical_marketing/repository"
	"gorm.io/gorm"
	"time"
)

const (
	TableName = "report"
)

type Model struct {
	ID int64 `bson:"id"`
	ReportUserID int64 `bson:"report_user_id"` // 报告人ID
	ConsumerMobile string `bson:"consumer_mobile"` // 客户手机号
	ConsumerName string `bson:"consumer_name"` // 客户姓名
	ExpectArriveTime *time.Time `bson:"except_arrive_time"` // 预期到访日期
	IsArrived int32 `bson:"is_arrived"` // 是否到访 1是0否
	ArrivedTime *time.Time `bson:"arrived_time"` // 到访时间
	PatientID int64 `bson:"patient_id"` // 患者ID
	CreateTime *time.Time `bson:"create_time"`
	IsDeleted int32 `bson:"is_deleted"` // 是否删除 1是 0 否
	DeleteTime *time.Time `bson:"delete_time"`
}

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
