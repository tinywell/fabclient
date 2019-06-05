package handler

import (
	"context"

	"github.com/tinywell/fabclient/pkg/sdk"
)

// Handler client handler for transaction
// type Handler interface {
// 	HandleMessage(ctx context.Context, msg Message)
// 	RegisterEvent() error
// 	GetEvent() <-chan sdk.Event
// 	RegisterTxHandleFunc(trancode TranCode, handleFunc TxHandleFunc) error
// 	RegisterSrcHandlerFunc(trancode TranCode, handleFunc SrcHandleFunc) error
// }

// TxHandleFunc transaction handle function
type TxHandleFunc func(ctx context.Context, msg Message, handler sdk.TxHandler)

// SrcHandleFunc resource management handle function
type SrcHandleFunc func(ctx context.Context, msg Message, handler sdk.ResourceManager)
