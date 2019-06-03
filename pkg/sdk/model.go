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
