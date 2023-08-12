package main

import (
	rebot "github.com/wangdzhao/community/discovery/rebot"
)

func startRebot() {
	robotName := "robot"
	bot := rebot.LoginWechat(robotName)
	if bot != nil {
		//rebot.WechatRobot(bot) //微信机器人
	}
}
