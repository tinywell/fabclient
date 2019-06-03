package common

// RspCode response code
type RspCode uint32

// response code
const (
	RspSuccess     = 200 // RspSuccess transaction success
	RspServerError = 500 // RspServerError internal server error
	RspTimeout     = 502 // RspTimeout transaction timeout ,check if block ok
)
