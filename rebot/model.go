package rebot

import "encoding/xml"

type Result struct {
	XMLName xml.Name `xml:"sysmsg"`
	MsgId   string   `xml:"revokemsg>msgid"`
}
