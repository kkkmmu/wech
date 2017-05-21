package message

import (
	"github.com/sirupsen/logrus"
)

/*

| MsgType | 说明 |
| ------- | --- |
| 1  | 文本消息 |
| 3  | 图片消息 |
| 34 | 语音消息 |
| 37 | 好友确认消息 |
| 40 | POSSIBLEFRIEND_MSG |
| 42 | 共享名片 |
| 43 | 视频消息 |
| 47 | 动画表情 |
| 48 | 位置消息 |
| 49 | 分享链接 |
| 50 | VOIPMSG |
| 51 | 微信初始化消息 |
| 52 | VOIPNOTIFY |
| 53 | VOIPINVITE |
| 62 | 小视频 |
| 9999 | SYSNOTICE |
| 10000 | 系统消息 |
| 10002 | 撤回消息 |

*/

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
