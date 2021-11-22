package messages

// Request to Server for calculate Fibonacci(Value)
type MsgRequest struct {
	Value int `json:"val"`
}

// response from Server with result calculation Fibonacci()
type MsgResponse struct {
	Value     int    `json:"val"`
	Result    int    `json:"res"`
	ErrorText string `json:"error"`
}
