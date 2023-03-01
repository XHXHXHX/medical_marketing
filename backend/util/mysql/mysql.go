package mysql

import (
	"errors"
	"fmt"
	"context"
	"strings"
	"time"

	"github.com/XHXHXHX/medical_marketing/util/logx"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

func NewMysql(host, user, password, dbName, logLevel string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	var dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s",
		user,
		password,
		host,
		dbName)
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         256,   // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  true,  // disable datetime precision support, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // smart configure based on used version
	}), &gorm.Config{
		Logger: &MysqlLogger{},
	})
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, errors.New("mysql db nil")
	}

	sqlDB, errSql := db.DB()
	if errSql != nil {
		return nil, err
	}
	// todo 开放参数 各个服务定制
	sqlDB.SetMaxIdleConns(8)
	sqlDB.SetMaxOpenConns(32)
	sqlDB.SetConnMaxLifetime(time.Minute * 8)

	if "debug" == logLevel {
		db = db.Debug()
	}

	db.Set("gorm:association", false)                //禁止自动创建/更新包含关系
	db.Set("gorm:association_save_reference", false) //禁止自动创建关联关系

	return db, nil
}

var (
	infoStr       = "%s\n[info] "
	warnStr       = "%s\n[warn] "
	errStr        = "%s\n[error] "
	traceStr      = "[FSQL] %s [%.3fms] [rows:%v] %s"
	traceWarnStr  = "[FSQL] %s %s [%.3fms] [rows:%v] %s"
	traceErrStr   = "[FSQL] %s %s [%.3fms] [rows:%v] %s"
	SlowThreshold = 200 * time.Millisecond
)

type MysqlLogger struct {}

func (l *MysqlLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *MysqlLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	logx.WithContext(ctx).Info(msg, args)
}

func (l *MysqlLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	logx.WithContext(ctx).Warn(msg, args)
}

func (l *MysqlLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	logx.WithContext(ctx).Error(msg, args)
}

func (l *MysqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound):
		sql, rows := fc()
		if rows == -1 {
			logx.WithContext(ctx).Error(traceErrStr, FileShortNameWithLineNum(utils.FileWithLineNum()), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logx.WithContext(ctx).Error(traceErrStr, FileShortNameWithLineNum(utils.FileWithLineNum()), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > SlowThreshold:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", SlowThreshold)
		if rows == -1 {
			logx.WithContext(ctx).Warn(traceWarnStr, FileShortNameWithLineNum(utils.FileWithLineNum()), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logx.WithContext(ctx).Warn(traceWarnStr, FileShortNameWithLineNum(utils.FileWithLineNum()), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	default:
		sql, rows := fc()
		if rows == -1 {
			logx.WithContext(ctx).Info(traceStr, FileShortNameWithLineNum(utils.FileWithLineNum()), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logx.WithContext(ctx).Info(traceStr, FileShortNameWithLineNum(utils.FileWithLineNum()), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func FileShortNameWithLineNum(fullPath string) string {
	paths := strings.Split(fullPath, "/")
	length := len(paths)
	if length > 0 {
		return paths[length-1]
	}
	return fullPath
}
