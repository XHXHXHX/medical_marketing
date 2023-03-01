package core

// Level 日志级别
type Level uint32

// 与logrus 日志级别对应
const (
	DebugLevel Level = iota + 2
	InfoLevel
	WarnLevel
	ErrorLevel
)

var LogLevel Level = DebugLevel

var (
	Namespace string
	Project   string
)

const skipOffset int = 2

// OffsetCtxKey _
type OffsetCtxKey struct{}

var (
	// SkipOffsetKey _
	SkipOffsetKey = OffsetCtxKey{}
)

const (
	msgSize     = 2048
	extDataSize = 2048
)
