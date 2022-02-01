package main

import (
	"fmt"

	"github.com/mailru/easyjson"

	"github.com/kyleconroy/wasm-json-greeter/hello"
)

func main() {
	fmt.Println("HELLO JSON")
	msg := hello.HelloRequest{Name: "foo"}
	blob, err := easyjson.Marshal(msg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := easyjson.Unmarshal(blob, &msg); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(msg.Name)
}
