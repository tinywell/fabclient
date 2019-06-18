package sdk

// Event fabric event
type Event struct {
	TxID        string
	CCName      string
	EventName   string
	BlockNumber uint64
	Payload     []byte
}

// RspMsg response message
type RspMsg struct {
	Code uint32
	TxID string
	Data []byte
}

// Request params for chaincode api
type Request struct {
	Channel      string
	Chaincode    string
	Fcn          string
	Args         [][]byte
	TransientMap map[string][]byte
}
