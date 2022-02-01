package hello

type HelloRequest struct {
	Name string `json:"name"`
}

// The response message containing the greetings
type HelloReply struct {
	Message string `json:"message"`
}
