package sdk

// TxHandler fabric transaction handler
type TxHandler interface {
	// Excute(channel, chaincode, fcn string, args [][]byte) RspMsg
	Excute(req Request) RspMsg
	// Query(channel, chaincode, fcn string, args [][]byte) RspMsg
	Query(req Request) RspMsg
	RegisterEvent(channel, ccName, event string) (<-chan *Event, error)
	UnRegisterEvent(channel, ccName, event string) error
}

// ResourceManager fabric resource manager,channel、chaincode etc.
type ResourceManager interface {
	CreateChannel()
	UpdateChannel()
	JoinChannel()
	InstallChaincode()
	InstantiateChaincode()
	UpgradeChaincode()
	LedgerQueryer
}

// LedgerQueryer fabric ledger queryer,block、chain info .etc
type LedgerQueryer interface {
	QueryBlock()
	QueryFirstBlock()
	QueryConfigBlock()
	QueryChainInfo()
}
