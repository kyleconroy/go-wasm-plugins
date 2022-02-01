greeter: hello.go main.go
	go build -o greeter .

linux-greeter: hello.go main.go
	GOOS=linux GOARCH=arm64 tinygo build -o linux-greeter -scheduler=none .
	
