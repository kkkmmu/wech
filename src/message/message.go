package message

import (
//	"github.com/sirupsen/logrus"
)

const (
	_              = iota
	TEXT           = 1     //| 文本消息 |
	IMAGE          = 3     //| 图片消息 |
	VOICE          = 34    //| 语音消息 |
	VERIFY         = 37    //| 好友确认消息 |
	POSSIBLEFRIEND = 40    //| POSSIBLEFRIEND_MSG |
	SHAREID        = 42    //| 共享名片 |
	VIDEO          = 43    //| 视频消息 |
	EMOJI          = 47    //| 动画表情 |
	SHAREPOSITION  = 48    //| 位置消息 |
	SHARELINK      = 49    //| 分享链接 |
	VOIP           = 50    //VOIPMSG |
	INIT           = 51    //| 微信初始化消息 |
	VOIPNOTIFY     = 52    //| VOIPNOTIFY |
	VOIPINVITE     = 53    //| VOIPINVITE |
	SHORTVIDEO     = 62    //| 小视频 |
	SYSNOTICE      = 9999  //| SYSNOTICE |
	SYSTEM         = 10000 // 系统消息 |
	REVOKE         = 10002 //| 撤回消息 |
)

var MsgTypeCount = 18
var MsgType = []int{
	TEXT,
	IMAGE,
	VOICE,
	VERIFY,
	POSSIBLEFRIEND,
	SHAREID,
	VIDEO,
	EMOJI,
	SHAREPOSITION,
	SHARELINK,
	VOIP,
	INIT,
	VOIPNOTIFY,
	VOIPINVITE,
	SHORTVIDEO,
	SYSNOTICE,
	SYSTEM,
	REVOKE,
}

var Away = "本人不在，请电话联系，谢谢"

type Message struct {
	MsgId                string        `json:"MsgId"`
	FromUserName         string        `json:"FromUserName"`
	ToUserName           string        `json:"ToUserName"`
	MsgType              int           `json:"MsgType"`
	Content              string        `json:"Content"`
	Status               int           `json:"Status"`
	ImgStatus            int           `json:"ImgStatus"`
	CreateTime           int           `json:"CreateTime"`
	VoiceLength          int           `json:"VoiceLength"`
	PlayLength           int           `json:"PlayLength"`
	FileName             string        `json:"FileName"`
	FileSize             string        `json:"FileSize"`
	MediaId              string        `json:"MediaId"`
	Url                  string        `json:"Url"`
	AppMsgType           int           `json:"AppMsgType"`
	StatusNotifyCode     int           `json:"StatusNotifyCode"`
	StatusNotifyUserName string        `json:"StatusNotifyUserName"`
	RecommendInfo        RecommendInfo `json:"RecommendInfo"`
	ForwardFlag          int           `json:"ForwardFlag"`
	AppInfo              AppInfo       `json:"AppInfo"`
	HasProductId         int           `json:"HasProductId"`
	Ticket               string        `json:"Ticket"`
	ImgHeight            int           `json:"ImgHeight"`
	ImgWidth             int           `json:"ImgWidth"`
	SubMsgType           int           `json:"SubMsgType"`
	NewMsgId             int64         `json:"NewMsgId"`
	OriContent           string        `json:"OriContent"`
}

func (msg *Message) IsTextMessage() bool {
	return msg.MsgType == TEXT
}

func (msg *Message) IsVoiceMessage() bool {
	return msg.MsgType == VOICE
}

func (msg *Message) IsVideoMessage() bool {
	return msg.MsgType == VIDEO
}

func (msg *Message) IsEmojiMessage() bool {
	return msg.MsgType == EMOJI
}

func (msg *Message) IsInitMessage() bool {
	return msg.MsgType == INIT
}

func (msg *Message) GetContent() string {
	return msg.Content
}

func (msg *Message) From() string {
	return msg.FromUserName
}

func (msg *Message) To() string {
	return msg.ToUserName
}

type TextMessage struct {
	Type         int
	Content      string
	FromUserName string
	ToUserName   string
	LocalID      int64
	ClientMsgId  int64
}

// MediaMessage
type MediaMessage struct {
	Type         int
	Content      string
	FromUserName string
	ToUserName   string
	LocalID      int64
	ClientMsgId  int64
	MediaId      string
}

// EmotionMessage: gif/emoji message struct
type EmotionMessage struct {
	ClientMsgId  int64
	EmojiFlag    int
	FromUserName string
	LocalID      int64
	MediaId      string
	ToUserName   string
	Type         int
}
