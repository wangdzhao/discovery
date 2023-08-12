package rebot

import (
	"fmt"

	util "github.com/wangdzhao/discovery/util"

	//"pinchequn/config"
	//"pinchequn/util"
	"regexp"
	"strings"

	"github.com/eatmoreapple/openwechat"
)

func GetGroupTypeMapInfo() map[string]string {
	groupTypeMap := map[string]string{
		//"文旅拼车群（10元每人）": config.GroupTypePincheMine,
		//"文旅10元拼车群（2群）": config.GroupTypePincheMine,
		//"天航10元拼车群":     config.GroupTypePincheMine,
		//"恒大文旅城回龙顺风车群":  config.GroupTypePincheOther,
		//"恒大文旅城拼车群（1）":  config.GroupTypePincheOther,
		//"黄龙溪谷清峯岸互助拼车群": config.GroupTypePincheOther,
		//"四方电脑客户服务群":    config.GroupTypeFatherGroup,
	}
	return groupTypeMap
}
func DeleteSlice(a []*openwechat.SentMessage, elem *openwechat.SentMessage) []*openwechat.SentMessage {
	tmp := make([]*openwechat.SentMessage, 0)
	for _, v := range a {
		if v != elem {
			tmp = append(tmp, v)
		}
	}
	return tmp
}
func GetSenderUserMap(sender *openwechat.User) map[string]string {
	userNickMap := make(map[string]string, 0)
	for _, v := range sender.MemberList {
		if v.DisplayName != "" {
			userNickMap[v.UserName] = v.DisplayName
			continue
		}
		userNickMap[v.UserName] = v.NickName

	}
	return userNickMap
}
func GetAllPincheGroupNames() []string {
	groups := GetGroupTypeMapInfo()
	groupNames := make([]string, 0, len(groups))
	//for name, v := range groups {
	//	//if v == config.GroupTypePincheMine || v == config.GroupTypePincheOther {
	//	//	groupNames = append(groupNames, name)
	//	//}
	//}
	return groupNames
}

func GetPincheGroupByType(gtype string) []string {
	groups := GetGroupTypeMapInfo()
	groupNames := make([]string, 0, len(groups))
	for name, v := range groups {
		if v == gtype {
			groupNames = append(groupNames, name)
		}
	}
	return groupNames
}
func FormatPincheInfo(msg *openwechat.Message) {
	reList := GetAllPincheGroupNames()
	repList := []string{"\n\n", "*"}
	if strings.Contains(msg.Content, "\n-") {
		sp := strings.Split(msg.Content, "\n-")
		msg.Content = sp[0]
	}
	for _, str := range reList {
		repList = append(repList, fmt.Sprintf("转自:%v,\n", str))
	}
	for _, str := range repList {
		if strings.Contains(msg.Content, str) {
			msg.Content = strings.ReplaceAll(msg.Content, str, "")
		}
	}
	msg.Content = strings.TrimRight(msg.Content, "\n")
	return
}

//func CheckPinCheGroup(groupUserName string, groupsTypeMap map[string]string) (string, bool) {
//	groupType := groupsTypeMap[groupUserName]
//	if groupType == config.GroupTypePincheMine || groupType == config.GroupTypePincheOther {
//		return groupType, true
//	}
//	return "", false
//}

func RemoveReplaceInfo(content string) string {
	groups := GetGroupTypeMapInfo()
	for key := range groups {
		replaceStr := fmt.Sprintf("转自:%v", key)
		if strings.Contains(content, replaceStr) {
			content = strings.ReplaceAll(content, replaceStr, "")
		}
	}
	content = cleanData(content)
	return content
}

//func CheckMsgIsWhite(msg string, senderNickeName string) bool {
//	whites := strings.Split(config.PinCheWhitelist, ",")
//	for _, name := range whites {
//		if name == senderNickeName {
//			return true
//		}
//	}
//	return false
//}

func CheckMsgIsPinCheInfo(msg string) bool {
	msgHaveDatas := []string{
		"【类型】",
		"【时间",
	}
	flag := true
	for _, str := range msgHaveDatas {
		if !strings.Contains(msg, str) {
			flag = false
			break
		}
	}
	return flag
}

func FormatGroupMessage(msg string) string {
	replacements := []struct {
		Old string
		New string
	}{
		{"【类型", "【类型】"},
		{"【时间", "【时间】"},
		{"【出发", "【出发】"},
		{"【备注", "【备注】"},
	}

	for _, r := range replacements {
		if strings.Contains(msg, r.Old) && !strings.Contains(msg, r.New) {
			msg = strings.ReplaceAll(msg, r.Old, r.New)
		}
	}

	return msg
}

func cleanData(data string) string {
	// 匹配需要清除的行数据
	re := regexp.MustCompile(`转自:.*\n`)
	// 将匹配到的行数据替换为空字符串
	cleanedData := re.ReplaceAllString(data, "")
	// 返回*之前的行数据和之后的行数据
	lines := strings.Split(cleanedData, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "*") {
			return strings.Join(lines[:i], "\n") + "\n\n********************\n" + strings.Join(lines[i:], "\n")
		}
	}
	return util.ProcessString(cleanedData)
}
