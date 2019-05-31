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
	RegisterTxHandleFunc(trancode string, handleFunc TxHandleFunc) error
	RegisterSrcHandlerFcun(trancode string, handleFunc SrcHandleFunc) error
	RegisterLedgerHandlerFcun(trancode string, handleFunc LedgerHandleFunc) error
}

// TxHandleFunc transaction handle function
type TxHandleFunc func(msg Message, handler sdk.TxHandler)

// SrcHandleFunc resource management handle function
type SrcHandleFunc func(msg Message, handler sdk.ResourceManager)

// LedgerHandleFunc ledger resource  management handle function
type LedgerHandleFunc func(msg Message, handler sdk.LedgerQueryer)
