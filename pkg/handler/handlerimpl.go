package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/tinywell/fabclient/pkg/common"
	"github.com/tinywell/fabclient/pkg/sdk"
)

var (
	defaultParams = Params{
		TxTimeout: 30,
		PoolSize:  200,
	}
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
	timeOut       time.Duration
	pool          chan struct{}
	// logger
}

// Core core parameters for handler Impl
type Core struct {
	TxHandler  sdk.TxHandler
	SrcManager sdk.ResourceManager
	// ...
}

// Params handler params
type Params struct {
	TxTimeout uint32 // Second
	PoolSize  uint32 // Concurrent handler size
}

// NewHandler return new Impl
func NewHandler(core Core, param Params) *Impl {
	if param.PoolSize == 0 {
		param.PoolSize = defaultParams.PoolSize
	}
	if param.TxTimeout == 0 {
		param.TxTimeout = defaultParams.TxTimeout
	}

	h := &Impl{
		msgs:          make(chan Message),
		events:        make(chan sdk.Event),
		txHandler:     core.TxHandler,
		txHandlers:    make(map[TranCode]TxHandleFunc),
		srcHandler:    core.SrcManager,
		srcHandlers:   make(map[TranCode]SrcHandleFunc),
		tranCodeStore: make(map[TranCode]TranType),
		doneC:         make(chan struct{}),
		timeOut:       time.Second * time.Duration(param.TxTimeout),
		pool:          make(chan struct{}, param.PoolSize),
	}
	go h.server()
	return h
}

func (h *Impl) server() {
	for {
		select {
		case msg := <-h.msgs:
			h.pool <- struct{}{}
			go h.handle(msg)
		case <-h.doneC:
			return
		}
	}
}

func (h *Impl) handle(msg Message) {
	defer func() {
		<-h.pool
	}()
	var rst Result
	trancode := msg.TranCode
	handlerType, ok := h.tranCodeStore[trancode]
	if !ok {
		// trancode not support,return error msg
		fmt.Println("trancode not support,return error msg")
		rst = Result{
			RspCode:  common.RspServerError,
			RspData:  []byte("trancode not support"),
			TranCode: msg.TranCode,
		}
		msg.Result <- rst
		return
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), h.timeOut)
	defer cancelFunc()
	switch handlerType {
	case TypeTxHandler:
		if h.txHandler == nil {
			h.doneC <- struct{}{}
			panic(fmt.Errorf("TxHandler not set"))
		}
		hfunc, ok := h.txHandlers[trancode]
		if !ok {
			// should not hanppen
			return
		}
		rst = hfunc(ctx, msg, h.txHandler)
	case TypeSrcHandler:
		if h.srcHandler == nil {
			h.doneC <- struct{}{}
			panic(fmt.Errorf("SrcManager not set"))
		}
		hfunc, ok := h.srcHandlers[trancode]
		if !ok {
			// should not hanppen
			return
		}
		rst = hfunc(ctx, msg, h.srcHandler)
	}
	msg.Result <- rst
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
	eventChan, err := h.txHandler.RegisterEvent(channel, ccName, event)
	if err != nil {
		return err
	}
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

// FillHandlerFunc regisnter handler functions in the box
func (h *Impl) FillHandlerFunc(box FuncBox) error {
	txHanlers := box.OpenTxHandlerBox()
	for k, v := range txHanlers {
		err := h.registerTxHandleFunc(k, v)
		if err != nil {
			return err
		}
	}
	srcHandlers := box.OpenSrcHandlerBox()
	for k, v := range srcHandlers {
		err := h.registerSrcHandlerFunc(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// RegisterTxHandleFunc register TxHandleFunc
func (h *Impl) registerTxHandleFunc(trancode TranCode, handleFunc TxHandleFunc) error {
	if _, ok := h.txHandlers[trancode]; ok {
		// already exist
		return fmt.Errorf("trancode %s already exist", trancode)
	}
	h.txHandlers[trancode] = handleFunc
	h.tranCodeStore[trancode] = TypeTxHandler
	return nil
}

// RegisterSrcHandlerFunc register SrcHandlerFunc
func (h *Impl) registerSrcHandlerFunc(trancode TranCode, handleFunc SrcHandleFunc) error {
	if _, ok := h.srcHandlers[trancode]; ok {
		// already exist
		return fmt.Errorf("trancode %s already exist", trancode)
	}
	h.srcHandlers[trancode] = handleFunc
	h.tranCodeStore[trancode] = TypeSrcHandler
	return nil
}
