package handler

// Message message for handler, result will be return throw channel Result
type Message struct {
	// ...
	TranCode string
	TranData []byte
	Result   chan<- Result
}

// Result handle result
type Result struct {
}
