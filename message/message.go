package message

import (
	"github.com/sirupsen/logrus"
)

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
