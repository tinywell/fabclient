package handler

type TranCode string

// Message message for handler, result will be return throw channel Result
type Message struct {
	// ...
	TranCode TranCode
	TranData []byte
	Result   chan<- Result
}

// Result handle result
type Result struct {
}

type HandlerType int

const (
	TypeTxHandler HandlerType = iota
	TypeEventHandler
	TypeLedgerhandler
)
