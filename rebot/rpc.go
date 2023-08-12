package rebot

import (
	"github.com/guonaihong/gout"
)

func ToGroupContent() (string, error) {
	contentGroup := ""
	err := gout.GET("http://127.0.0.1:8009/news/gropu_day").BindBody(&contentGroup).Do()
	if err != nil {
		return "", err
	}
	return contentGroup, nil
}

func ToFaContent() (string, error) {
	contentFa := ""
	err := gout.GET("http://127.0.0.1:8009/news/today").BindBody(&contentFa).Do()
	if err != nil {
		return "", err
	}
	return contentFa, nil
}

//func sendToWechatMsgUrl(body string) {
//	gout.POST(config.WechatMsgUrl).SetBody(&body).Do()
//	time.Sleep(time.Microsecond * 200)
//}
