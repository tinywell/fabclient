package handler

// TranCode transaction code
type TranCode string

// RspCode response code
type RspCode uint32

// response code
const (
	RspSuccess     RspCode = 200 // RspSuccess transaction success
	RspServerError RspCode = 500 // RspServerError internal server error
	RspTimeout     RspCode = 502 // RspTimeout transaction timeout ,check if block ok
)

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
	RepData  []byte
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
