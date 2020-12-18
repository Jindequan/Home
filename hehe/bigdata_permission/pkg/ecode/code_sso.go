package ecode

//sso错误码格式，五位数，10开头。eg：10xxx

var (
	UnregisteredCamcard = New(10001) //该后缀邮箱未注册CC（目前只有bertadata.com需要注册cc）
	UnsupportedEmail    = New(10002) //不支持该后缀邮箱
	ForceUpdatePassword = New(10003) //密码长时间未修改
	PlatformNotExist    = New(10004) //平台不存在
)

var SsoCodeMsg = map[int]string{
	UnregisteredCamcard.Code(): "该后缀邮箱未注册CC",
	UnsupportedEmail.Code():    "不支持该后缀邮箱",
	ForceUpdatePassword.Code(): "密码长时间未修改，请先修改密码",
	PlatformNotExist.Code():    "平台不存在",
}
