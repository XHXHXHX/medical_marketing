package errs

type Error struct {
	Code string
	Message string
}

var (
	ExistSameConsumerMobile = NewSimpleError("report.existsSameMobile", "存在相同客户手机号")
)