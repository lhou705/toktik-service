package common

const (
	Success            = 0  // 请求成功
	NotLogin           = -1 // 未登录
	LoginStatusExpired = -2 // 登录状态过期
	InvalidFiledType   = -3 // 字段类型错误
	FiledIsNone        = -4 // 字段为空
	FiledLengthInvalid = -5
	ReqError           = -6 // 请求失败
	InvalidFiledValue  = -7 // 字段值错误
)
