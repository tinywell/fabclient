package handler

import (
	"context"

	core "github.com/tinywell/fabclient/pkg/handler"
	"github.com/tinywell/fabclient/pkg/sdk"
)

// ExpHandler example handler function
type ExpHandler struct {
	channel   string
	chaincode string
	user      string
}

// CreateHandler return new example handler
func CreateHandler(channel, chaincode, user string) *ExpHandler {
	return &ExpHandler{
		channel:   channel,
		chaincode: chaincode,
		user:      user,
	}
}

// SaveData handler save data,trancode is 'EXP100'
func (h *ExpHandler) SaveData(ctx context.Context, msg core.Message, handler sdk.TxHandler) {
	rspmsg := handler.Excute(h.channel, h.chaincode, "savedata", [][]byte{msg.TranData})
	rst := core.Result{
		RspCode:  rspmsg.Code,
		RspData:  rspmsg.Data,
		TranCode: msg.TranCode,
		TxID:     rspmsg.TxID,
	}
	msg.Result <- rst
}

// ReadData handler read data,trancode is 'EXP200'
func (h *ExpHandler) ReadData(ctx context.Context, msg core.Message, handler sdk.TxHandler) {
	rspmsg := handler.Excute(h.channel, h.chaincode, "readdata", [][]byte{msg.TranData})
	rst := core.Result{
		RspCode:  rspmsg.Code,
		RspData:  rspmsg.Data,
		TranCode: msg.TranCode,
		TxID:     rspmsg.TxID,
	}
	msg.Result <- rst
}
