package customer_task

import (
	"github.com/XHXHXHX/medical_marketing/util/common"
	"time"
)

const (
	Finished IsFinished = 1
	UnFinish IsFinished = 2
)

type IsFinished int64

func (f IsFinished) IsValid() bool {
	return f == Finished || f == UnFinish
}

func (f IsFinished) IsFinished() bool {
	return f == Finished
}

type CustomerTask struct {
	ID int64 `bson:"id"`
	ReportID int64 `bson:"report_id"`
	UserID int64 `bson:"user_id"`
	DistributeTime time.Time `bson:"distribute_time"`
	IsFinished IsFinished
	FinishTime *time.Time
}

type CustomerTaskHistory struct {
	ID int64 `bson:"id"`
	TaskID int64 `bson:"task_id"`
	Desc string `bson:"desc"`
	CreateTime time.Time `bson:"create_time"`
	UserID int64 `bson:"user_id"`
}

type SelectListRequest struct {
	UserIDs []int64
	ReportIDs []int64
	IsFinished IsFinished
	DistributeStartTime *time.Time
	DistributeEndTime *time.Time
	FinishStartTime *time.Time
	FinishEndTime *time.Time
	Page *common.Page
}