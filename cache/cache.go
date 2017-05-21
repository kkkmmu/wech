package cache

import (
	"message"
	"net/http"
	"synckey"
)

type Cache struct {
	BaseResponse *message.BaseResponse `json:"BaseResponse"`
	DeviceID     string                `json:"DeviceID"`
	SyncKey      *synckey.SyncKey      `json:"SyncKey"`
	UserName     string                `json:"UserName"`
	NickName     string                `json:"NickName"`
	SyncServer   string                `json:"SyncServer"`
	Ticket       string                `json:"Ticket"`
	Lang         string                `json:"Lang"`
	Cookies      []*http.Cookie        `json:"Cookies"`
}
