package sms

type Sms interface {
	Send(phoneNo string) bool
}

type AliYun struct {
}

func (aliYun *AliYun) Send(phoneNo string) bool {
	return true
}
