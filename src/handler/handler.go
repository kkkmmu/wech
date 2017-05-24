package handler

import (
	"errors"
	"log"
	"message"
)

var (
	InvalidMsgTypeError = errors.New("Invalid message Type")
)

type HandleFunc func(msg *message.LocalMessage) *message.LocalMessage

type Handler struct {
	name    string
	state   bool
	msgtype int
	handle  HandleFunc
}

func NewHandler(name string, msgtype int, state bool, handle HandleFunc) (*Handler, error) {
	if !message.IsValid(msgtype) {
		return nil, InvalidMsgTypeError
	}

	return &Handler{
		name:    name,
		state:   state,
		msgtype: msgtype,
		handle:  handle,
	}, nil
}

func (h *Handler) Enable() {
	h.state = true
}

func (h *Handler) Disable() {
	h.state = false
}

func (h *Handler) Type() int {
	return h.msgtype
}

func (h *Handler) Name() string {
	return h.name
}

func (h *Handler) IsEnabled() bool {
	return h.state
}

func (h *Handler) Handle(msg *message.LocalMessage) *message.LocalMessage {
	if h.state {
		return h.handle(msg)
	}
	return nil
}

func (h *Handler) SetName(name string) *Handler {
	h.name = name
	return h
}

func (h *Handler) SetState(state bool) *Handler {
	h.state = state
	return h
}

func (h *Handler) SetType(t int) *Handler {
	h.msgtype = t
	return h
}

func (h *Handler) SetHandleFunc(handle HandleFunc) *Handler {
	h.handle = handle
	return h
}

func init() {
	DefaultHandlers = make([]*Handler, 0, len(message.Type))
	for t, n := range message.Type {
		h, err := NewHandler(n, t, true, defaultMessageHandler)
		if err != nil {
			log.Println("Create default message handler failed for: ", n, t, err)
			continue
		}

		DefaultHandlers = append(DefaultHandlers, h)
	}
}
