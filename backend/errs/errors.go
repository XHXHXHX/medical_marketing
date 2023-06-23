package errs

var (
	InvalidParams = NewSimpleError("common.InvalidParams", "无效参数")
	NotFoundData = NewSimpleError("common.NotFoundData", "无法找到相应数据")

	ExistSameConsumerMobile = NewSimpleError("report.existsSameMobile", "存在相同客户手机号")
	BanRecover = NewSimpleError("report.BanRecover", "回访的报单不能撤回")
	ExpectBeforeNow = NewSimpleError("report.ExpectBeforeNow", "预期到访时间不能早于当前时间")

	NotBelongWithYou = NewSimpleError("report.NotBelongWithYou", "当前客服不属于你")

	UserPasswordError = NewSimpleError("user.UserPasswordError", "密码错误")
	RoleError = NewSimpleError("user.RoleError", "直属领导或超级管理员才可以办理离职")
	MobileRepeat = NewSimpleError("user.MobileRepeat", "手机号重复")

	ExpireToken = NewSimpleError("auth.ExpireToken", "认证令牌过期")
	InvalidToken = NewSimpleError("auth.InvalidToken", "认证令牌无效")
	NoRole = NewSimpleError("auth.NoRole", "无权限")
)