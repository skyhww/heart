package cfg

type AliYunConfig struct {
	AccessKeyId     string     `yaml:"access_key_id"`
	AccessKeySecret string     `yaml:"access_key_secret"`
	RegionId        string     `yaml:"region_id"`
	SmsConfig       *SmsConfig `yaml:"sms"`
}

type SmsConfig struct {
	SendSmsConfig *SendSmsConfig `yaml:"send_sms_code"`
}
type SendSmsConfig struct {
	Code         string `yaml:"code"`
	Method       string `yaml:"method"`
	Scheme       string `yaml:"scheme"`
	SignName     string `yaml:"sign_name"`
	TemplateCode string `yaml:"template_code"`
	Domain       string `yaml:"domain"`
	Version      string `yaml:"version"`
	ApiName      string `yaml:"api_name"`
}

