package logrus

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/XHXHXHX/medical_marketing/util/logx/core"
)

func NewProvider() *core.Provider {
	log := logrus.New()
	log.Formatter = &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	}
	log.SetLevel(logrus.DebugLevel)

	return core.NewProvider(func(ts int64, level core.Level, msg string, args map[string]interface{}) {
		entry := log.WithFields(args)
		switch level {
		case core.DebugLevel:
			{
				entry.Debug(msg)
			}
		case core.InfoLevel:
			{
				entry.Info(msg)
			}
		case core.WarnLevel:
			{
				entry.Warn(msg)
			}
		case core.ErrorLevel:
			{
				entry.Error(msg)
			}
		default:
			{
				//do nothing
			}
		}
	})
}
