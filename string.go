package main

//https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage?ticket=AQkJglrmmmXwfjPT39C1FJht@qrticket_0&uuid=oeApj3RRLQ==&lang=zh_CN&scan=1494147471

import (
	"log"
	"regexp"
	"strings"
)

var str = "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage?ticket=AQkJglrmmmXwfjPT39C1FJht@qrticket_0&uuid=oeApj3RRLQ==&lang=zh_CN&scan=1494147471"
var fields = "ticket=AQkJglrmmmXwfjPT39C1FJht@qrticket_0&uuid=oeApj3RRLQ==&lang=zh_CN&scan=1494147471"

var re = regexp.MustCompile(`ticket=(?P<ticket>[[:word:]_\$@#\?\-=]+)&uuid=(?P<uuid>[[:word:]_\$@#\?\-=]+)&lang=(?P<lang>[[:word:]_\$@#\?\-=]+)&scan=(?P<lang>[[:word:]_\$@#\?\-=]+)`)

func main() {
	strs := strings.Split(str, "?")
	log.Println(strs[1])
	m := re.FindStringSubmatch(fields)
	log.Println(m)
	for _, v := range m {
		log.Println(v)
	}
}
