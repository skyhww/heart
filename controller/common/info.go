package common

import "heart/service"

var PassWordLengthNotEnough = &service.Info{Code: "000001", Message: "密码长度必须大于6"}
var RequestDataRequired = &service.Info{Code: "000002", Message: "上送的数据为空"}
var PasswordRequired = &service.Info{Code: "000003", Message: "密码不能为空"}
var IllegalMobileFormat = &service.Info{Code: "000004", Message: "手机号非法"}
var ConfirmPasswordNotMatched = &service.Info{Code: "000005", Message: "确认密码不匹配"}
var SmsCodeRequired = &service.Info{Code: "000006", Message: "短信验证码不能为空"}
var IllegalSmsCodeFormat = &service.Info{Code: "000008", Message: "短信验证码格式不正确"}
var MobileRequired = &service.Info{Code: "000009", Message: "短信验证码格式不正确"}
