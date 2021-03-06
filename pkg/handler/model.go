package handler

// TranCode transaction code
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
	RspCode  uint32
	RspData  []byte
	TranCode TranCode
	TxID     string
	Err      error
}

// TranType handler function type
type TranType int

// handler func type
const (
	TypeTxHandler TranType = iota
	TypeSrcHandler
)
