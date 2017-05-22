package wechat

import (
	"log"
	"message"
)

type Handler func(wc *WeChat, in *message.Message)

type Plugin struct {
	name    string
	msgType int
	enable  bool
	handle  Handler
}

func NewPlugin(name string, msgType int, handler Handler) *Plugin {
	return &Plugin{
		name:    name,
		enable:  false,
		handle:  handler,
		msgType: msgType,
	}
}

var DefaultPlugins []*Plugin

func (p *Plugin) Name() string {
	return p.name
}

func (p *Plugin) Enable() {
	p.enable = true
}

func (p *Plugin) Disable() {
	p.enable = false
}

func (p *Plugin) IsEnabled() bool {
	return p.enable
}

func (p *Plugin) Handle(wc *WeChat, msg *message.Message) {
	p.handle(wc, msg)
}

func defaultTextMessageHandler(wc *WeChat, msg *message.Message) {
	log.Println("Received TEXT message from: ", wc.ContactDBUserName[msg.FromUserName].NickName, " content: ", msg.Content)
	wc.SendTextMessage(wc.UserName, msg.FromUserName, message.Away)
}

func defaultImageMessageHandler(wc *WeChat, msg *message.Message) {
	log.Println("Received IMAGE message from: ", wc.ContactDBUserName[msg.FromUserName].NickName, " content: ", msg.Content)
	wc.SendTextMessage(wc.UserName, msg.FromUserName, message.Away)
}

func defaultVoiceMessageHandler(wc *WeChat, msg *message.Message) {
	log.Println("Received VOICE message from: ", wc.ContactDBUserName[msg.FromUserName].NickName, " content: ", msg.Content)
	wc.SendTextMessage(wc.UserName, msg.FromUserName, message.Away)
}

func defaultVerifyMessageHandler(wc *WeChat, msg *message.Message) {
	log.Println("Received VERIFY message from: ", wc.ContactDBUserName[msg.FromUserName].NickName, " content: ", msg.Content)
	wc.SendTextMessage(wc.UserName, msg.FromUserName, message.Away)
}

func defaultPossibleFriendMessageHandler(wc *WeChat, msg *message.Message) {
	log.Println("Received POSSIBLEFRIEND message from: ", wc.ContactDBUserName[msg.FromUserName].NickName, " content: ", msg.Content)
	wc.SendTextMessage(wc.UserName, msg.FromUserName, message.Away)
}

func defaultShareIDMessageHandler(wc *WeChat, msg *message.Message) {
	log.Println("Received SHAREID message from: ", wc.ContactDBUserName[msg.FromUserName].NickName, " content: ", msg.Content)
	wc.SendTextMessage(wc.UserName, msg.FromUserName, message.Away)
}

func defaultVideoMessageHandler(wc *WeChat, msg *message.Message) {
	log.Println("Received VIDEO message from: ", wc.ContactDBUserName[msg.FromUserName].NickName, " content: ", msg.Content)
	wc.SendTextMessage(wc.UserName, msg.FromUserName, message.Away)
}

func defaultEmojiMessageHandler(wc *WeChat, msg *message.Message) {
	log.Println("Received EMOJI message from: ", wc.ContactDBUserName[msg.FromUserName].NickName, " content: ", msg.Content)
	wc.SendTextMessage(wc.UserName, msg.FromUserName, message.Away)
}

func defaultSharePositionMessageHandler(wc *WeChat, msg *message.Message) {
	log.Println("Received SHAREPOSITION message from: ", wc.ContactDBUserName[msg.FromUserName].NickName, " content: ", msg.Content)
	wc.SendTextMessage(wc.UserName, msg.FromUserName, message.Away)
}

func defaultShareLinkMessageHandler(wc *WeChat, msg *message.Message) {
	log.Println("Received SHARELINK message from: ", wc.ContactDBUserName[msg.FromUserName].NickName, " content: ", msg.Content)
	wc.SendTextMessage(wc.UserName, msg.FromUserName, message.Away)
}

func defaultVOIPMessageHandler(wc *WeChat, msg *message.Message) {
	log.Println("Received VOIP message from: ", wc.ContactDBUserName[msg.FromUserName].NickName, " content: ", msg.Content)
	wc.SendTextMessage(wc.UserName, msg.FromUserName, message.Away)
}

func defaultInitMessageHandler(wc *WeChat, msg *message.Message) {
	log.Println("Received INIT message from: ", wc.ContactDBUserName[msg.FromUserName].NickName, " content: ", msg.Content)
	wc.SendTextMessage(wc.UserName, msg.FromUserName, message.Away)
}

func defaultVOIPNotifyMessageHandler(wc *WeChat, msg *message.Message) {
	log.Println("Received VOIPNOTIFY message from: ", wc.ContactDBUserName[msg.FromUserName].NickName, " content: ", msg.Content)
	wc.SendTextMessage(wc.UserName, msg.FromUserName, message.Away)
}

func defaultVOIPInviteMessageHandler(wc *WeChat, msg *message.Message) {
	log.Println("Received VOIPINVITE message from: ", wc.ContactDBUserName[msg.FromUserName].NickName, " content: ", msg.Content)
	wc.SendTextMessage(wc.UserName, msg.FromUserName, message.Away)
}

func defaultShortVideoMessageHandler(wc *WeChat, msg *message.Message) {
	log.Println("Received SHORTVIDEO message from: ", wc.ContactDBUserName[msg.FromUserName].NickName, " content: ", msg.Content)
	wc.SendTextMessage(wc.UserName, msg.FromUserName, message.Away)
}

func defaultSystemNoticeMessageHandler(wc *WeChat, msg *message.Message) {
	log.Println("Received SYSTEMNOTICE message from: ", wc.ContactDBUserName[msg.FromUserName].NickName, " content: ", msg.Content)
	wc.SendTextMessage(wc.UserName, msg.FromUserName, message.Away)
}

func defaultSystemMessageHandler(wc *WeChat, msg *message.Message) {
	log.Println("Received SYSTEM message from: ", wc.ContactDBUserName[msg.FromUserName].NickName, " content: ", msg.Content)
	wc.SendTextMessage(wc.UserName, msg.FromUserName, message.Away)
}

func defaultRevokeMessageHandler(wc *WeChat, msg *message.Message) {
	log.Println("Received INVOKE message from: ", wc.ContactDBUserName[msg.FromUserName].NickName, " content: ", msg.Content)
	wc.SendTextMessage(wc.UserName, msg.FromUserName, message.Away)
}
