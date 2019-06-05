package sdk

// MockTxHandler mock implement for TxHandler
type MockTxHandler struct {
}

// Excute mock chaincode execution
func (mock *MockTxHandler) Excute(channel, chaincode, fcn string, args [][]byte) RspMsg {
	return RspMsg{
		Code: 200,
		TxID: "TESTMOCK0001",
		Data: []byte("Hello World"),
	}
}

// Query mock chaincode query
func (mock *MockTxHandler) Query(channel, chaincode, fcn string, args [][]byte) RspMsg {
	return RspMsg{
		Code: 200,
		TxID: "TESTMOCK0002",
		Data: []byte("Hello World"),
	}
}

// RegisterEvent mock register chaincode event
func (mock *MockTxHandler) RegisterEvent(channel, ccName, event string) (<-chan *Event, error) {
	return nil, nil
}

// UnRegisterEvent mock register chaincode event
func (mock *MockTxHandler) UnRegisterEvent(channel, ccName, event string) error {
	return nil
}

// mock ResourceManager
func (mock *MockTxHandler) CreateChannel()        {}
func (mock *MockTxHandler) UpdateChannel()        {}
func (mock *MockTxHandler) JoinChannel()          {}
func (mock *MockTxHandler) InstallChaincode()     {}
func (mock *MockTxHandler) InstantiateChaincode() {}
func (mock *MockTxHandler) UpgradeChaincode()     {}

// mock LedgerQueryer
func (mock *MockTxHandler) QueryBlock()       {}
func (mock *MockTxHandler) QueryFirstBlock()  {}
func (mock *MockTxHandler) QueryConfigBlock() {}
func (mock *MockTxHandler) QueryChainInfo()   {}
