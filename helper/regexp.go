package helper

import "regexp"

var MobileRegexp = regexp.MustCompile(`^1([38][0-9]|14[57]|5[^4])\d{8}$`)
