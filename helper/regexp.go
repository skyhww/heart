package helper

import "regexp"

var MobileRegexp = regexp.MustCompile(`^1([38][0-9]|14[57]|5[^4])\d{8}$`)
var SmsCodeRegexp = regexp.MustCompile(`^\d{6}$`)
//最少6位，包括至少1个大写字母，1个小写字母，1个数字，1个特殊字符
//var PasswordRegexp=regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[$@$!%*?&])[A-Za-z\d$@$!%*?&]{6,}$/`)
//4到16位（字母，数字，下划线，减号）
var UsernameRegexp = regexp.MustCompile(`/^[a-zA-Z]\w{3,15}$/`)
