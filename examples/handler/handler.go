package handler

import (
	"context"
	"encoding/json"

	"github.com/tinywell/fabclient/pkg/common"

	core "github.com/tinywell/fabclient/pkg/handler"
	"github.com/tinywell/fabclient/pkg/sdk"
)

// ExpHandler example handler function
type ExpHandler struct {
	channel   string
	chaincode string
	user      string
}

// PutKey param for putkey
type PutKey struct {
	Name    string `json:"name,omitempty"`
	KeyCert string `json:"key_cert,omitempty"`
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
func (h *ExpHandler) SaveData(ctx context.Context, msg core.Message, handler sdk.TxHandler) core.Result {
	putkey := &PutKey{}
	err := json.Unmarshal(msg.TranData, putkey)
	if err != nil {
		return core.Result{
			RspCode:  common.RspServerError,
			RspData:  []byte(err.Error()),
			TranCode: msg.TranCode,
		}
	}
	rspmsg := handler.Excute(h.channel, h.chaincode, "keyadd", [][]byte{[]byte(putkey.Name), []byte(putkey.KeyCert)})
	return core.Result{
		RspCode:  rspmsg.Code,
		RspData:  rspmsg.Data,
		TranCode: msg.TranCode,
		TxID:     rspmsg.TxID,
	}
}

// ReadData handler read data,trancode is 'EXP200'
func (h *ExpHandler) ReadData(ctx context.Context, msg core.Message, handler sdk.TxHandler) core.Result {
	rspmsg := handler.Query(h.channel, h.chaincode, "keyget", [][]byte{msg.TranData})
	return core.Result{
		RspCode:  rspmsg.Code,
		RspData:  rspmsg.Data,
		TranCode: msg.TranCode,
		TxID:     rspmsg.TxID,
	}
}

// OpenTxHandlerBox return txhandler functions
func (h *ExpHandler) OpenTxHandlerBox() map[core.TranCode]core.TxHandleFunc {
	boxMap := make(map[core.TranCode]core.TxHandleFunc)
	boxMap[core.TranCode("EXP100")] = h.SaveData
	boxMap[core.TranCode("EXP200")] = h.ReadData
	return boxMap
}

// OpenSrcHandlerBox return src handler functions
func (h *ExpHandler) OpenSrcHandlerBox() map[core.TranCode]core.SrcHandleFunc {
	boxMap := make(map[core.TranCode]core.SrcHandleFunc)
	return boxMap
}
