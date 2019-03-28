package sms

type Sms interface {
	Send(phoneNo string,content string) bool
}

type AliYun struct {
}

func (aliYun *AliYun) Send(phoneNo string,content string) bool {
	/*client, err := sdk.NewClientWithAccessKey("default", "<accessKeyId>", "<accessSecret>")
	if err != nil {
		panic(err)
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "QuerySendDetails"
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Print(response.GetHttpContentString())*/
	return true
}

