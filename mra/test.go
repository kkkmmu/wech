package main

import (
	"wechat"
)

func main() {
	wc, err := wechat.NewWeChatClient("test")
	if err != nil {
		panic(err)
	}

	if err := wc.FastLogin(); err != nil {
		wc.GetUUID()
		wc.GetQRCode()
		go wc.WaitForQRCodeScan()
		<-wc.Login
	}
	wc.GetBaseRequest()
	wc.WeChatInit()
	wc.StatusNotify()
	wc.GetContactList()
	wc.GetGroupMemberList()
	// wc.GetSyncServer()
	wc.Run()
}
