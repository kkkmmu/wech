package wechat

import (
	"message"
	"strings"
)

var adminPlugin = &Plugin{
	name:    "admin",
	msgType: message.TEXT,
	enable:  true,
	handle:  pluginManagerMessageHandler,
}

func pluginManagerMessageHandler(wc *WeChat, msg *message.Message) {
	if msg.IsTextMessage() {
		content := msg.GetContent()
		content = strings.TrimSpace(content)
		if strings.HasPrefix(content, "#admin") {
			fields := strings.Split(content, " ")
			command := strings.TrimSpace(fields[1])
			switch command {
			case "list":
				wc.ShowPluginList(msg.FromUserName)
			case "help":
				wc.ShowHelp(msg.FromUserName)
			case "enable":
				wc.EnablePluginByName(msg.FromUserName, strings.TrimSpace(fields[2]))
			case "disable":
				wc.DisablePluginByName(msg.FromUserName, strings.TrimSpace(fields[2]))
			default:
				wc.SendTextMessage(wc.UserName, msg.FromUserName, "未知命令，请重新输入")
			}
		}
	}
}

func (wc *WeChat) ShowHelp(to string) {

}

func (wc *WeChat) ShowPluginList(to string) {

}

func (wc *WeChat) EnablePluginByName(to, name string) {

}

func (wc *WeChat) DisablePluginByName(to, name string) {

}
