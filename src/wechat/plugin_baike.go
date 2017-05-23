package wechat

import (
	"io/ioutil"
	"log"
	"message"
	"net/http"
	"strings"
)

var BaiKePlugin = &Plugin{
	name:    "Baike",
	msgType: message.TEXT,
	enable:  true,
	handle:  BaikeHandler,
}

func BaikeHandler(wc *WeChat, msg *message.Message) {
	if msg.IsTextMessage() {
		data := msg.GetContent()
		if strings.HasPrefix(data, "?bk") {
			data = data[3:len(data)]
			resp, err := http.Get("http://baike.baidu.com/item/" + data)
			if err != nil {
				log.Println("Cannot get: ", data, " from baike with error: ", err.Error())
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Println(err)
					return
				}

				log.Println(string(body))
				wc.SendShareLinkMessage(wc.UserName, msg.FromUserName, "http://baike.baidu.com/item/"+data)
				return
			}

			log.Println(resp.Location())
			log.Println(resp.Status)
			log.Println(resp.StatusCode)
			location, _ := resp.Location()
			if resp.StatusCode == 302 {
				content, err := http.Get("http://baike.baidu.com" + location.String())
				if err != nil {
					log.Println("Cannot redirect for: ", data, " with error: ", err.Error())
					return
				}
				defer content.Body.Close()

				body, err := ioutil.ReadAll(content.Body)
				if err != nil {
					log.Println("Cannot get content")
					return
				}

				wc.SendShareLinkMessage(wc.UserName, msg.FromUserName, "http://baike.baidu.com/item/"+location.String())
				log.Println(string(body))
			}
		}
	}
}
