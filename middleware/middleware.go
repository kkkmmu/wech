package middleware

import (
	"message"
)

type MessageHandler func(message *message.Message) error

type Middleware struct {
	name       string
	enable     bool
	msgHandler MessageHandler
}

func (mw *Middleware) Name() string {
	return mw.name
}

func (mw *Middleware) Enable() {
	mw.enable = true
}

func (mw *Middleware) Disable() {
	mw.enable = false
}

func (mw *Middleware) HandleMessage(msg *message.Message) error {
	if mw.enable {
		return mw.HandleMessage(msg)
	}
	return nil
}
