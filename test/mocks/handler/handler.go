// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/handler/handler.go

// Package mhandler is a generated GoMock package.
package mhandler

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	handler "github.com/tinywell/fabclient/pkg/handler"
	sdk "github.com/tinywell/fabclient/pkg/sdk"
	reflect "reflect"
)

// MockHandler is a mock of Handler interface
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// HandleMessage mocks base method
func (m *MockHandler) HandleMessage(ctx context.Context, msg handler.Message) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleMessage", ctx, msg)
}

// HandleMessage indicates an expected call of HandleMessage
func (mr *MockHandlerMockRecorder) HandleMessage(ctx, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleMessage", reflect.TypeOf((*MockHandler)(nil).HandleMessage), ctx, msg)
}

// RegisterEvent mocks base method
func (m *MockHandler) RegisterEvent() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterEvent")
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterEvent indicates an expected call of RegisterEvent
func (mr *MockHandlerMockRecorder) RegisterEvent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterEvent", reflect.TypeOf((*MockHandler)(nil).RegisterEvent))
}

// GetEvent mocks base method
func (m *MockHandler) GetEvent() <-chan sdk.Event {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvent")
	ret0, _ := ret[0].(<-chan sdk.Event)
	return ret0
}

// GetEvent indicates an expected call of GetEvent
func (mr *MockHandlerMockRecorder) GetEvent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvent", reflect.TypeOf((*MockHandler)(nil).GetEvent))
}

// RegisterTxHandleFunc mocks base method
func (m *MockHandler) RegisterTxHandleFunc(trancode handler.TranCode, handleFunc handler.TxHandleFunc) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterTxHandleFunc", trancode, handleFunc)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterTxHandleFunc indicates an expected call of RegisterTxHandleFunc
func (mr *MockHandlerMockRecorder) RegisterTxHandleFunc(trancode, handleFunc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterTxHandleFunc", reflect.TypeOf((*MockHandler)(nil).RegisterTxHandleFunc), trancode, handleFunc)
}

// RegisterSrcHandlerFunc mocks base method
func (m *MockHandler) RegisterSrcHandlerFunc(trancode handler.TranCode, handleFunc handler.SrcHandleFunc) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterSrcHandlerFunc", trancode, handleFunc)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterSrcHandlerFunc indicates an expected call of RegisterSrcHandlerFunc
func (mr *MockHandlerMockRecorder) RegisterSrcHandlerFunc(trancode, handleFunc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterSrcHandlerFunc", reflect.TypeOf((*MockHandler)(nil).RegisterSrcHandlerFunc), trancode, handleFunc)
}
