package sms

type Sms interface {
	Send(phoneNo string,content string) bool
}

type AliYun struct {
}

func (aliYun *AliYun) Send(phoneNo string,content string) bool {
	return true
}
