package handler

import (
	"context"

	"github.com/tinywell/fabclient/pkg/sdk"
)

// Handler client handler for transaction
type Handler interface {
	HandleMessage(ctx context.Context, msg Message)
	RegisterEvent() error
	GetEvent() <-chan sdk.Event
	// RegisterTxHandleFunc(trancode TranCode, handleFunc TxHandleFunc) error
	// RegisterSrcHandlerFunc(trancode TranCode, handleFunc SrcHandleFunc) error
	FillHandlerFunc(box FuncBox) error
}

// TxHandleFunc transaction handle function
type TxHandleFunc func(ctx context.Context, msg Message, handler sdk.TxHandler) Result

// SrcHandleFunc resource management handle function
type SrcHandleFunc func(ctx context.Context, msg Message, handler sdk.ResourceManager) Result

// FuncBox handler functions with TranCode
type FuncBox interface {
	OpenTxHandlerBox() map[TranCode]TxHandleFunc
	OpenSrcHandlerBox() map[TranCode]SrcHandleFunc
}
