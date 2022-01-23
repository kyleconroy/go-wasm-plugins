.PHONY: clean

build: example plugin.wasm

example: hello/hello.pb.go main.go
	go build -o example

hello/hello.pb.go: hello/hello.proto
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--go-vtproto_out=. \
		--go-vtproto_opt=paths=source_relative,features=marshal+unmarshal+size \
		hello/hello.proto

plugin.wasm: hello/hello.pb.go main.go
	tinygo build -o plugin.wasm -wasm-abi=generic -target=wasi

clean:
	rm -f tinyexample
	rm -f example
	rm -f hello/hello.pb.go
