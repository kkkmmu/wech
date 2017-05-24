package handler

import (
	"message"
)

var DefaultHandlers []*Handler

var DefaultResponse = &message.LocalMessage{}

func defaultMessageHandler(msg *message.LocalMessage) *message.LocalMessage {
	return DefaultResponse.SetFrom(msg.To()).SetTo(msg.From()).SetType(message.TEXT).SetContent([]byte("本人微信不在线，请电话联系，谢谢！"))
}
