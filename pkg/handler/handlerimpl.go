package handler

import (
	"context"

	"github.com/tinywell/fabclient/pkg/sdk"
)

// Impl implement for handler
type Impl struct {
	msgs          chan Message
	events        chan sdk.Event
	txHandler     sdk.TxHandler
	txHandlers    map[string]TxHandleFunc
	tranCodeStore map[string]HandlerType
	doneC         chan struct{}
	// logger
}

// Core core parameters for handler Impl
type Core struct {
	TxHandler sdk.TxHandler
	// ...
}

// NewHandler return new Impl
func NewHandler(core Core) (*Impl, error) {
	h := &Impl{
		msgs:          make(chan Message, 100),
		txHandler:     core.TxHandler,
		txHandlers:    make(map[string]TxHandleFunc),
		tranCodeStore: make(map[string]HandlerType),
		doneC:         make(chan struct{}),
	}
	go h.server()
	return h, nil
}

func (h *Impl) server() {
	for {
		select {
		case _ = <-h.msgs:
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

	}
}

// RegisterEvent register fabric  event,event messger will be boradcast throw channel h.events
func (h *Impl) RegisterEvent() error {
	return nil
}

// GetEvent return event channel
func (h *Impl) GetEvent() <-chan sdk.Event {
	return h.events
}

// RegisterTxHandleFunc register TxHandleFunc
func (h *Impl) RegisterTxHandleFunc(trancode string, handleFunc TxHandleFunc) error {
	return nil
}

// RegisterSrcHandlerFunc register SrcHandlerFunc
func (h *Impl) RegisterSrcHandlerFunc(trancode string, handleFunc SrcHandleFunc) error {
	return nil
}

// RegisterLedgerHandlerFunc register  LedgerHandlerFunc
func (h *Impl) RegisterLedgerHandlerFunc(trancode string, handleFunc LedgerHandleFunc) error {
	return nil
}
