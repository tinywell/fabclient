package sdk

// TxHandler fabric transaction handler
type TxHandler interface {
	Excute()
	Query()
}

// EventHandler fabric event handler
type EventHandler interface {
	RegisterEvent()
}

// ResourceManager fabric resource manager,channel、chaincode etc.
type ResourceManager interface {
	CreateChannel()
	UpdateChannel()
	JoinChannel()
	InstallChaincode()
	InstantiateChaincode()
	UpgradeChaincode()
}

// LedgerQueryer fabric ledger queryer,block、chain info .etc
type LedgerQueryer interface {
	QueryBlock()
	QueryFirstBlock()
	QueryConfigBlock()
	QueryChainInfo()
}
