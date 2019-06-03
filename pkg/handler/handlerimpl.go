package handler

import (
	"context"
	"fmt"

	"github.com/tinywell/fabclient/pkg/sdk"
)

// Impl implement for handler
type Impl struct {
	msgs          chan Message
	events        chan sdk.Event
	txHandler     sdk.TxHandler
	txHandlers    map[TranCode]TxHandleFunc
	srcHandler    sdk.ResourceManager
	srcHandlers   map[TranCode]SrcHandleFunc
	tranCodeStore map[TranCode]TranType
	doneC         chan struct{}
	// logger
}

// Core core parameters for handler Impl
type Core struct {
	TxHandler  sdk.TxHandler
	SrcManager sdk.ResourceManager
	// ...
}

// NewHandler return new Impl
func NewHandler(core Core) (*Impl, error) {
	h := &Impl{
		msgs:          make(chan Message, 100),
		txHandler:     core.TxHandler,
		txHandlers:    make(map[TranCode]TxHandleFunc),
		srcHandler:    core.SrcManager,
		srcHandlers:   make(map[TranCode]SrcHandleFunc),
		tranCodeStore: make(map[TranCode]TranType),
		doneC:         make(chan struct{}),
	}
	go h.server()
	return h, nil
}

func (h *Impl) server() {
	for {
		select {
		case msg := <-h.msgs:
			trancode := msg.TranCode
			handlerType, ok := h.tranCodeStore[trancode]
			if !ok {
				// trancode not support,return error msg
				continue
			}

			switch handlerType {
			case TypeTxHandler:
				if h.txHandler == nil {
					panic(fmt.Errorf("TxHandler not set"))
				}
				hfunc, ok := h.txHandlers[trancode]
				if !ok {
					// should not hanppen

				}
				hfunc(msg, h.txHandler)
			case TypeSrcHandler:
				if h.txHandler == nil {
					panic(fmt.Errorf("SrcManager not set"))
				}
				hfunc, ok := h.srcHandlers[trancode]
				if !ok {
					// should not hanppen

				}
				hfunc(msg, h.srcHandler)
			}
		case <-h.doneC:
			return
		}
	}
}

// HandleMessage handle message
func (h *Impl) HandleMessage(ctx context.Context, msg Message) {
	select {
	case h.msgs <- msg:
	case <-ctx.Done():
		return
	}
}

// RegisterEvent register fabric event, event messges will be boradcast throw channel h.events
func (h *Impl) RegisterEvent(channel, ccName string, event string) error {
	eventChan := h.txHandler.RegisterEvent(channel, ccName, event)
	go func(event <-chan *sdk.Event) {
		for e := range event {
			h.events <- *e
		}
	}(eventChan)
	return nil
}

// GetEvent return event channel
func (h *Impl) GetEvent() <-chan sdk.Event {
	return h.events
}

// RegisterTxHandleFunc register TxHandleFunc
func (h *Impl) RegisterTxHandleFunc(trancode TranCode, handleFunc TxHandleFunc) error {
	if _, ok := h.txHandlers[trancode]; ok {
		// already exist
		return nil
	}
	h.txHandlers[trancode] = handleFunc
	return nil
}

// RegisterSrcHandlerFunc register SrcHandlerFunc
func (h *Impl) RegisterSrcHandlerFunc(trancode TranCode, handleFunc SrcHandleFunc) error {
	if _, ok := h.srcHandlers[trancode]; ok {
		// already exist
		return nil
	}
	h.srcHandlers[trancode] = handleFunc
	return nil
}
