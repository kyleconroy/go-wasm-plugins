package main

import (
	"fmt"

	"github.com/kyleconroy/wasm-greeter/hello"
)

func main() {
	fmt.Println("HELLO WORLD")
	msg := hello.HelloRequest{Name: "foo"}
	return

	blob, err := msg.MarshalVT()
	if err != nil {
			fmt.Println(err.Error())
			return
	}
	if err := msg.UnmarshalVT(blob); err != nil {
			fmt.Println(err.Error())
			return
	}
	fmt.Println(msg.Name)
}
