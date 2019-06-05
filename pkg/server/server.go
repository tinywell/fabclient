package server

import (
	"github.com/tinywell/fabclient/pkg/handler"
	"github.com/tinywell/fabclient/pkg/sdk"
)

// Server client server
type Server interface {
	ReceiveMessage() <-chan handler.Message
	HandleEvent(event sdk.Event)
}
