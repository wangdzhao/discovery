package rebot

import (
	"fmt"
	"io"
	"time"

	util "github.com/wangdzhao/discovery/util"

	"github.com/eatmoreapple/openwechat"
)

func LoginWechat(robotName string) *openwechat.Bot {
	bot := openwechat.DefaultBot(openwechat.Desktop)
	bot.UUIDCallback = generateQrCodeCallback

	reloadStorage := openwechat.NewFileHotReloadStorage(fmt.Sprint(robotName, ".json"))
	defer func(reloadStorage io.ReadWriteCloser) {
		_ = reloadStorage.Close()
	}(reloadStorage)

	// 执行热登录
	if err := bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption()); err != nil {
		util.LogError("登录失败: %v", err)
		return nil
	}

	return bot
}

func generateQrCodeCallback(uuid string) {
	qrcodeUrl := openwechat.GetQrcodeUrl(uuid)
	util.SaveCodeImg(qrcodeUrl)
	codeMd5, codeBase, _ := util.SendImageToWeChatRobot("code.png")
	msgData := "登录地址：" + qrcodeUrl + "\n" + "创建时间：" + time.Now().Format("2006-01-02 15:04:05")

	imageBody := `{
		"msgtype": "image",
		"image": {
			"base64": "%s",
			"md5": "%s"
		}
	}`

	textBody := `{
		"msgtype": "text",
		"text": {
			"content": "%v"
		}
	}`

	imageBody = fmt.Sprintf(imageBody, codeBase, codeMd5)
	textBody = fmt.Sprintf(textBody, msgData)
	//sendToWechatMsgUrl(textBody)
	//sendToWechatMsgUrl(imageBody)
}

//func PinCheGroupMsg(sender *openwechat.User, msg *openwechat.Message, groups *openwechat.Groups, groupsTypeMap map[string]string) []*openwechat.SentMessage {
//	groupUserName := getGroupUserName(sender, msg)
//
//	groupType, typeOk := CheckPinCheGroup(groupUserName, groupsTypeMap)
//	if !typeOk {
//		return nil
//	}
//
//	if isMessageDuplicate(msg) {
//		return nil
//	}
//
//	userMap := GetSenderUserMap(sender)
//	senderNickeName := getSenderNickname(msg, userMap)
//	FormatPincheInfo(msg)
//	isWhite := CheckMsgIsWhite(msg.Content, senderNickeName)
//	checkRes := CheckMsgIsPinCheInfo(msg.Content)
//
//	if groupType == config.GroupTypePincheMine && !isWhite && !checkRes {
//		sendInvalidFormatReply(senderNickeName, msg)
//		return nil
//	}
//
//	if !checkRes {
//		return nil
//	}
//
//	sentList := ForwardPinCheMsg(sender, groups, msg)
//	storeMessage(msg)
//	return sentList
//}
//
//func getGroupUserName(sender *openwechat.User, msg *openwechat.Message) string {
//	if msg.IsSendBySelf() {
//		return msg.ToUserName
//	}
//	return sender.UserName
//}
//
//func isMessageDuplicate(msg *openwechat.Message) bool {
//	key := util.Sha256(msg.Content)
//	redis := goredis.NewRedisHelper()
//	val, _ := redis.Get(redis.Context(), key).Result()
//	return val == "true"
//}
//
//func getSenderNickname(msg *openwechat.Message, userMap map[string]string) string {
//	senderInGroupUserNameData := strings.Split(msg.RawContent, ":<br/>")
//	return userMap[senderInGroupUserNameData[0]]
//}
//
//func sendInvalidFormatReply(senderNickeName string, msg *openwechat.Message) {
//	msg.ReplyText(fmt.Sprintf("@%v 非常抱歉，您发布的信息格式有误，请按照下面的模板发布信息", senderNickeName))
//	time.Sleep(time.Microsecond * 1000)
//	msg.ReplyText("【类型】车找人/人找车\n【时间】\n【出发】\n【到达】\n【车位】\n【备注】")
//}
//
//func storeMessage(msg *openwechat.Message) {
//	key := util.Sha256(msg.Content)
//	redis := goredis.NewRedisHelper()
//	redis.Set(redis.Context(), key, "true", config.ReplayMsg*time.Minute)
//	redis.Set(redis.Context(), "msg_"+msg.MsgId, msg.Content, 2*time.Minute)
//}
//func WechatRobot(bot *openwechat.Bot) {
//	groupTypeMap := GetGroupTypeMapInfo()
//	self, err := bot.GetCurrentUser()
//	if err != nil {
//		util.LogError("Failed to get current user: %v", err)
//		return
//	}
//
//	firends, err := self.Friends()
//	if err != nil {
//		util.LogError("Failed to get friends: %v", err)
//		return
//	}
//
//	groups, err := self.Groups()
//	if err != nil {
//		util.LogError("Failed to get group: %vs", err)
//		return
//	}
//
//	groupNameTypeMap, redisG := processGroups(groups, groupTypeMap, firends)
//
//	RedisMsg(firends, redisG)
//	DayNoticeMsg(firends, redisG)
//
//	sentList := []*openwechat.SentMessage{}
//	bot.MessageHandler = func(msg *openwechat.Message) {
//		if msg.IsText() && msg.IsSendByGroup() {
//			sender, err := msg.Sender()
//			if err != nil {
//				util.LogError("Failed to get sender: %v", err)
//				return
//			}
//			sentListItem := PinCheGroupMsg(sender, msg, &groups, groupNameTypeMap)
//			sentList = append(sentList, sentListItem...)
//		}
//	}
//}
//
//func processGroups(groups []*openwechat.Group, groupTypeMap map[string]string, firends openwechat.Friends) (map[string]string, []*openwechat.Group) {
//	groupNameTypeMap := make(map[string]string, 0)
//	redisG := []*openwechat.Group{}
//	for _, g := range groups {
//		groupType := groupTypeMap[g.NickName]
//		switch groupType {
//		case config.GroupTypePincheMine, config.GroupTypePincheOther:
//			groupNameTypeMap[g.UserName] = groupType
//			redisG = append(redisG, g)
//		case config.GroupTypeFatherGroup:
//			FatherTask(firends, g)
//		}
//	}
//	return groupNameTypeMap, redisG
//}
//
//func ForwardPinCheMsg(sender *openwechat.User, groups *openwechat.Groups, msg *openwechat.Message) []*openwechat.SentMessage {
//	groupNames := GetAllPincheGroupNames()
//	nickName := sender.NickName
//
//	if msg.IsSendBySelf() {
//		group := groups.GetByUsername(msg.ToUserName)
//		nickName = group.NickName
//	}
//
//	sentList := []*openwechat.SentMessage{}
//	groupMsg := FormatGroupMessage(msg.Content)
//
//	for _, gname := range groupNames {
//		if gname == nickName {
//			continue
//		}
//
//		group := groups.GetByNickName(gname)
//		if group == nil {
//			util.LogError(fmt.Sprintf("Failed to get group by nickname: %s", gname), nil)
//			continue
//		}
//
//		groupMsg = RemoveReplaceInfo(groupMsg)
//		sent, err := group.SendText(fmt.Sprintf("%v \n\n********************\n%v", groupMsg, fmt.Sprintf("转自:%v", nickName)))
//		if err != nil {
//			util.LogError("Failed to send text message: %v", err)
//			continue
//		}
//
//		sent.SendMessage.MediaId = msg.MsgId
//		sentList = append(sentList, sent)
//		time.Sleep(time.Microsecond * 50)
//	}
//
//	return sentList
//}
