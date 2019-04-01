package common

import (
	"heart/service/common"
)

var PassWordLengthNotEnough = &base.Info{Code: "000001", Message: "密码长度必须大于6"}
var RequestDataRequired = &base.Info{Code: "000002", Message: "上送的数据为空"}
var PasswordRequired = &base.Info{Code: "000003", Message: "密码不能为空"}
var IllegalMobileFormat = &base.Info{Code: "000004", Message: "手机号非法"}
var ConfirmPasswordNotMatched = &base.Info{Code: "000005", Message: "确认密码不匹配"}
var SmsCodeRequired = &base.Info{Code: "000006", Message: "短信验证码不能为空"}
var IllegalSmsCodeFormat = &base.Info{Code: "000008", Message: "短信验证码格式不正确"}
var MobileRequired = &base.Info{Code: "000009", Message: "短信验证码格式不正确"}
var UploadFailed = &base.Info{Code: "000010", Message: "上传失败"}
var ReLogin = &base.Info{Code: "000011", Message: "请重新登录"}
var IllegalRequest = &base.Info{Code: "000012", Message: "读取数据异常"}
var IllegalRequestDataFormat = &base.Info{Code: "000013", Message: "数据格式异常"}
