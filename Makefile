.PHONY: run clean

GOMODCACHE := $(shell go env GOMODCACHE)

run: cli/cli
	cd cli && ./cli

clean:
	rm -f cli/*.wasm
	rm -f plugin_go/plugin.wasm
	rm -f plugin_go_json/plugin.wasm
	rm -f plugin_rust/plugin.wasm

cli/plugin_rust.wasm: plugin_rust/src/main.rs plugin_rust/Cargo.toml
	cd plugin_rust && cargo build --target wasm32-wasi
	cp plugin_rust/target/wasm32-wasi/debug/plugin.wasm cli/plugin_rust.wasm
	
# cli/plugin.wasm: plugin/src/main.rs plugin/Cargo.toml
# 	# Overwrite the built plugin with the one from tinygo
# 	cp ../go-wasm-plugin/go_plugin/plugin.wasm cli/plugin.wasm
#
plugin_go_json/hello/hello_easyjson.go: plugin_go_json/hello/hello.go
	cd plugin_go_json/hello && easyjson -all hello.go

plugin_go_json/plugin.wasm: plugin_go_json/hello/hello_easyjson.go plugin_go_json/main.go
	cd plugin_go_json && tinygo build -o plugin.wasm -gc=leaking -scheduler=none -wasm-abi=generic -target=wasi

cli/plugin_go_json.wasm: plugin_go_json/plugin.wasm
	cp plugin_go_json/plugin.wasm cli/plugin_go_json.wasm

plugin_go/hello/hello.pb.go: plugin_go/hello/hello.proto
	cd plugin_go && protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		--go-vtproto_out=. \
		--go-vtproto_opt=paths=source_relative,features=marshal+unmarshal+size \
		hello/hello.proto

plugin_go/plugin.wasm: plugin_go/hello/hello.pb.go plugin_go/main.go
	cd plugin_go && tinygo build -o plugin.wasm -gc=leaking -scheduler=none -wasm-abi=generic -target=wasi

cli/plugin_go.wasm: plugin_go/plugin.wasm
	cp plugin_go/plugin.wasm cli/plugin_go.wasm

cli/cli: cli/plugin_rust.wasm cli/plugin_go.wasm cli/plugin_go_json.wasm cli/main.go ${GOMODCACHE}/github.com/bytecodealliance/wasmtime-go@v0.33.0/build/macos-aarch64/libwasmtime.a
	cd cli && go build .

${GOMODCACHE}/github.com/bytecodealliance/wasmtime-go@v0.33.0/build/macos-aarch64/libwasmtime.a: ${GOMODCACHE}/github.com/bytecodealliance/wasmtime-go@v0.33.0/build/macos-aarch64 ./wasmtime/target/release
	sudo cp wasmtime/target/release/libwasmtime.a ${GOMODCACHE}/github.com/bytecodealliance/wasmtime-go@v0.33.0/build/macos-aarch64

${GOMODCACHE}/github.com/bytecodealliance/wasmtime-go@v0.33.0/build/macos-aarch64:
	sudo mkdir -p ${GOMODCACHE}/github.com/bytecodealliance/wasmtime-go@v0.33.0/build/macos-aarch64

wasmtime/target/release: wasmtime
	cd wasmtime && cargo build -p wasmtime-c-api --release

wasmtime:
	git clone https://github.com/bytecodealliance/wasmtime.git --depth 1
	cd wasmtime && git fetch --tags
	cd wasmtime && git checkout v0.33.0
	cd wasmtime && git submodule update --init
