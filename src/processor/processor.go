package processor

import (
	"errors"
	"fmt"
	"handler"
	"log"
	"message"
	"strings"
	"sync"
)

var (
	AlreadExistError        = errors.New("Save handler Alread exist")
	NotExistError           = errors.New("Handler does not exist")
	InvalidHandlerTypeError = errors.New("Invalid message Type")
	InvalidMessageTypeError = errors.New("Invalid message Type")
)

var DefaultProcessor *Processor

var adminPrefix = "%am"

var commandList = "%am show\n%am help\n%am enable all\n%am enable xxx\n%am disable all\n%am disable xxx\n"

type Processor struct {
	handlerDB  map[int][]*handler.Handler
	allHandler map[string]*handler.Handler
	Lock       *sync.Mutex
}

func NewProcessor() *Processor {
	return &Processor{
		handlerDB:  make(map[int][]*handler.Handler, len(message.Type)),
		allHandler: make(map[string]*handler.Handler, len(message.Type)),
		Lock:       &sync.Mutex{},
	}
}

//Register A new message handler
func (p *Processor) RegisterHandler(h *handler.Handler) error {
	if _, ok := p.allHandler[h.Name()]; ok {
		return AlreadExistError
	}

	p.Lock.Lock()
	defer p.Lock.Unlock()
	if _, ok := p.handlerDB[h.Type()]; !ok {
		p.handlerDB[h.Type()] = make([]*handler.Handler, 0, 1)
	}

	p.allHandler[h.Name()] = h
	p.handlerDB[h.Type()] = append(p.handlerDB[h.Type()], h)

	return nil
}

func (p *Processor) DeRegisterHandler(h *handler.Handler) error {
	if _, ok := p.allHandler[h.Name()]; !ok {
		return NotExistError
	}

	p.Lock.Lock()
	defer p.Lock.Unlock()

	for i, d := range p.handlerDB[h.Type()] {
		if d.Name() == h.Name() {
			delete(p.allHandler, h.Name())
			p.handlerDB[h.Type()] = append(p.handlerDB[h.Type()][:i], p.handlerDB[h.Type()][i+1:]...)
			break
		}
	}

	return nil
}

func (p *Processor) Process(msg *message.LocalMessage) ([]*message.LocalMessage, error) {
	if !message.IsValid(msg.Type()) {
		log.Println("Received message of unknown type: ", msg.Type())
		return nil, InvalidMessageTypeError
	}

	log.Println("Received ", message.Type[msg.Type()], " message from: ", msg.From(), " content: ", string(msg.Content()))
	if IsAdminMessage(msg) {
		return []*message.LocalMessage{p.processAdminMessage(msg)}, nil
	}

	return p.processNormalMessage(msg), nil
}

func (p *Processor) processNormalMessage(msg *message.LocalMessage) []*message.LocalMessage {
	response := make([]*message.LocalMessage, 0, 1)
	p.Lock.Lock()
	defer p.Lock.Unlock()
	for _, h := range p.handlerDB[msg.Type()] {
		response = append(response, h.Handle(msg))
	}

	return response
}

func IsAdminMessage(msg *message.LocalMessage) bool {
	content := string(msg.Content())
	content = strings.TrimSpace(content)

	if strings.HasPrefix(content, adminPrefix) {
		return true
	}

	return false
}

func (p *Processor) processAdminMessage(msg *message.LocalMessage) *message.LocalMessage {
	content := string(msg.Content())
	fields := strings.Split(content, " ")

	log.Println("Received Management Message: ", content)
	if fields[1] == "show" {
		return p.DumpCurrentStatus(msg)
	} else if fields[1] == "help" {
		return p.processShowCommandList(msg)
	} else if len(fields) == 3 {
		if fields[1] == "enable" {
			if fields[2] == "all" {
				return p.processEnableAllHandler(msg)
			}

			msg.SetContent([]byte(fields[2]))
			return p.processEnableHandler(msg)

		} else if fields[1] == "disable" {
			if fields[2] == "all" {
				return p.processDisableAllHandler(msg)
			}
			msg.SetContent([]byte(fields[2]))
			return p.processDisableHandler(msg)
		}
	}

	return p.DumpCurrentStatus(msg)
}

func (p *Processor) GetAllHandler() []string {
	p.Lock.Lock()
	defer p.Lock.Unlock()
	handlers := make([]string, 0, len(p.allHandler))
	for name, _ := range p.allHandler {
		handlers = append(handlers, name)
	}
	return handlers
}

func (p *Processor) DisableHandlerByName(name string) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()
	h, ok := p.allHandler[name]
	if !ok {
		return NotExistError
	}

	h.Disable()

	return nil
}

func (p *Processor) EnableHandlerByName(name string) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()
	h, ok := p.allHandler[name]
	if !ok {
		return NotExistError
	}

	h.Disable()
	return nil
}

func (p *Processor) DisableHandlerByType(t int) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	if !message.IsValid(t) {
		return InvalidHandlerTypeError
	}

	for _, h := range p.handlerDB[t] {
		h.Disable()
	}

	return nil
}

func (p *Processor) EnableHandlerByType(t int) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	if !message.IsValid(t) {
		return InvalidHandlerTypeError
	}

	for _, h := range p.handlerDB[t] {
		h.Enable()
	}

	return nil
}

func (p *Processor) DisableAllHandler() {
	for _, h := range p.allHandler {
		h.Disable()
	}
}

func (p *Processor) EnableAllHandler() {
	for _, h := range p.allHandler {
		h.Enable()
	}
}

func (p *Processor) IsEnabled(name string) bool {
	if h, ok := p.allHandler[name]; ok {
		return h.IsEnabled()
	}
	return false
}

func (p *Processor) DumpCurrentStatus(msg *message.LocalMessage) *message.LocalMessage {
	handlers := p.GetAllHandler()

	content := fmt.Sprintf("Total handler count: %d\n", len(handlers))
	for _, h := range handlers {
		content = content + "[" + h + "]" + ":"
		if p.IsEnabled(h) {
			content = content + "enable\n"
		} else {
			content = content + "disable\n"
		}
	}

	response := &message.LocalMessage{}
	return response.SetFrom(msg.To()).SetTo(msg.From()).SetType(message.TEXT).SetContent([]byte(content))
}

func (p *Processor) processShowCommandList(msg *message.LocalMessage) *message.LocalMessage {
	response := &message.LocalMessage{}
	return response.SetFrom(msg.To()).SetTo(msg.From()).SetType(message.TEXT).SetContent([]byte(commandList))
}

func (p *Processor) processEnableHandler(msg *message.LocalMessage) *message.LocalMessage {
	p.EnableHandlerByName(string(msg.Content()))
	return p.DumpCurrentStatus(msg)
}

func (p *Processor) processDisableHandler(msg *message.LocalMessage) *message.LocalMessage {
	p.DisableHandlerByName(string(msg.Content()))
	return p.DumpCurrentStatus(msg)
}

func (p *Processor) processDisableAllHandler(msg *message.LocalMessage) *message.LocalMessage {
	p.DisableAllHandler()
	return p.DumpCurrentStatus(msg)
}

func (p *Processor) processEnableAllHandler(msg *message.LocalMessage) *message.LocalMessage {
	p.EnableAllHandler()
	return p.DumpCurrentStatus(msg)
}

func init() {
	DefaultProcessor = NewProcessor()
	log.Println(len(message.Type))
	log.Println(len(DefaultProcessor.allHandler))
	log.Println(DefaultProcessor.allHandler)
	log.Println(len(handler.DefaultHandlers))
	for _, h := range handler.DefaultHandlers {
		DefaultProcessor.RegisterHandler(h)
	}
}
