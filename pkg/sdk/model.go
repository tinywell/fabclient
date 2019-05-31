package sdk

// Event fabric event
type Event struct {
	TxID        string
	CCName      string
	EventName   string
	BlockNumber uint64
	Payload     []byte
}
