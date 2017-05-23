package plugin

import (
	"message"
)

type MessageHandler func(msg *message.Message) (resp *message.Message)

type Plugin struct {
	name    string
	msgType int
	handle  MessageHandler
}

func NewPlugin(name string, msgType int, handler MessageHandler) (*Plugin, error) {
	//TODO: Prameter Check
	return &Plugin{
		Name:    name,
		msgType: msgType,
		handle:  handler,
	}, nil
}

func (p *Plugin) Handle(msg *message.Message) *message.Message {

}
