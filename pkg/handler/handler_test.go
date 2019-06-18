package handler

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/tinywell/fabclient/pkg/sdk"
	msdk "github.com/tinywell/fabclient/test/mocks/sdk"
)

var (
	testChannel = "testchannel"
	testCC      = "testcc"
)

// func TestRegisterTxHandleFunc(t *testing.T) {
// 	h := createHandler(t)
// 	txHF := func(ctx context.Context, msg Message, handler sdk.TxHandler) Result {
// 		return Result{}
// 	}
// 	err := h.RegisterTxHandleFunc("TEST001", txHF)
// 	if err != nil {
// 		t.Errorf("RegisterTxHandleFunc err: %v", err)
// 	}
// }

// func TestRegisterSrcHandlerFunc(t *testing.T) {
// 	h := createHandler(t)
// 	srcHF := func(ctx context.Context, msg Message, handler sdk.ResourceManager) Result {
// 		return Result{}
// 	}
// 	err := h.RegisterSrcHandlerFunc("TEST101", srcHF)
// 	if err != nil {
// 		t.Errorf("RegisterSrcHandlerFunc err: %v", err)
// 	}
// }

func TestGetEvent(t *testing.T) {
	h := createHandler(t)
	eventC := h.GetEvent()
	if eventC == nil {
		t.Error("get event error")
	}
}

func TestFillBox(t *testing.T) {
	h := createHandler(t)
	err := h.FillHandlerFunc(testBox{})
	if err != nil {
		t.Error(err)
	}
}

func TestHandleMessage(t *testing.T) {
	ctl := gomock.NewController(t)
	txHandler := msdk.NewMockTxHandler(ctl)
	rsp1 := sdk.RspMsg{
		Code: 200,
		Data: []byte("hello programer"),
		TxID: "TXIDTEST0001",
	}
	
	txHandler.EXPECT().Excute(sdk.Request{Channel: testChannel, Chaincode: testCC, Fcn: "savedata", Args: [][]byte{[]byte("hello world")}}).Return(rsp1)
	rsp2 := sdk.RspMsg{
		Code: 200,
		Data: []byte("hello world"),
		TxID: "TXIDTEST1001",
	}
	
	txHandler.EXPECT().Query(sdk.Request{Channel: testChannel, Chaincode: testCC, Fcn: "readdata", Args: [][]byte{[]byte("TXIDTEST0001")}}).Return(rsp2)
	srcManager := msdk.NewMockResourceManager(ctl)
	core := Core{
		TxHandler:  txHandler,
		SrcManager: srcManager,
	}
	params := Params{
		TxTimeout: 10,
		PoolSize:  100,
	}
	h := NewHandler(core, params)
	mockBox := testBox{}
	err := h.FillHandlerFunc(mockBox)
	if err != nil {
		t.Error(err)
	}

	rstC1 := make(chan Result)
	msg1 := Message{
		TranCode: "TEST002",
		TranData: []byte("hello world"),
		Result:   rstC1,
	}

	rstC2 := make(chan Result)
	msg2 := Message{
		TranCode: "TEST102",
		TranData: []byte("TXIDTEST0001"),
		Result:   rstC2,
	}
	ctx := context.Background()
	ctxTimeout1, cancelFunc1 := context.WithTimeout(ctx, time.Second*30)
	defer cancelFunc1()
	h.HandleMessage(ctx, msg1)
	select {
	case rst := <-rstC1:
		if rst.RspCode != 200 {
			t.Error("handler execute message error")
		}
		t.Log(rst)
	case <-ctxTimeout1.Done():
		t.Error("handle message timeout")
	}
	h.HandleMessage(ctx, msg2)
	ctxTimeout2, cancelFunc2 := context.WithTimeout(ctx, time.Second*30)
	defer cancelFunc2()
	select {
	case rst := <-rstC2:
		if rst.RspCode != 200 {
			t.Error("handler execute message error")
		}
		t.Log(rst)
	case <-ctxTimeout2.Done():
		t.Error("handle message timeout")
	}
}

func createHandler(t *testing.T) *Impl {
	ctl := gomock.NewController(t)
	txHandler := msdk.NewMockTxHandler(ctl)
	// txHandler.EXPECT().Excute("", "", "", nil).Return(nil)
	srcManager := msdk.NewMockResourceManager(ctl)
	core := Core{
		TxHandler:  txHandler,
		SrcManager: srcManager,
	}
	params := Params{
		TxTimeout: 10,
		PoolSize:  100,
	}
	h := NewHandler(core, params)
	return h
}

type testBox struct {
}

func (box testBox) OpenTxHandlerBox() map[TranCode]TxHandleFunc {
	boxMap := make(map[TranCode]TxHandleFunc)
	txHF1 := func(ctx context.Context, msg Message, handler sdk.TxHandler) Result {
		rsp := handler.Excute(sdk.Request{Channel: testChannel, Chaincode: testCC, Fcn: "savedata", Args: [][]byte{msg.TranData}})
		return Result{
			RspCode:  200,
			RspData:  rsp.Data,
			TranCode: msg.TranCode,
			TxID:     rsp.TxID,
		}
	}
	txHF2 := func(ctx context.Context, msg Message, handler sdk.TxHandler) Result {		
		rsp := handler.Query(sdk.Request{Channel: testChannel, Chaincode: testCC, Fcn: "readdata", Args: [][]byte{msg.TranData}})
		return Result{
			RspCode:  200,
			RspData:  rsp.Data,
			TranCode: msg.TranCode,
			TxID:     rsp.TxID,
		}
	}
	boxMap[TranCode("TEST002")] = txHF1
	boxMap[TranCode("TEST102")] = txHF2
	return boxMap
}

func (box testBox) OpenSrcHandlerBox() map[TranCode]SrcHandleFunc {
	boxMap := make(map[TranCode]SrcHandleFunc)
	return boxMap
}
