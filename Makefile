.PHONY: run

run: cli/plugin.wasm cli/cli
	cd cli && ./cli

cli/plugin.wasm: plugin/src/main.rs plugin/Cargo.toml
	cd plugin && cargo build --target wasm32-wasi
	cp plugin/target/wasm32-wasi/debug/plugin.wasm cli/plugin.wasm

cli/cli: cli/main.go /Users/kyle/go/pkg/mod/github.com/bytecodealliance/wasmtime-go@v0.33.0/build/macos-aarch64/libwasmtime.a
	cd cli && go build .

/Users/kyle/go/pkg/mod/github.com/bytecodealliance/wasmtime-go@v0.33.0/build/macos-aarch64/libwasmtime.a: /Users/kyle/go/pkg/mod/github.com/bytecodealliance/wasmtime-go@v0.33.0/build/macos-aarch64 ./wasmtime/target/release
	sudo cp wasmtime/target/release/libwasmtime.a /Users/kyle/go/pkg/mod/github.com/bytecodealliance/wasmtime-go@v0.33.0/build/macos-aarch64

/Users/kyle/go/pkg/mod/github.com/bytecodealliance/wasmtime-go@v0.33.0/build/macos-aarch64:
	sudo mkdir -p /Users/kyle/go/pkg/mod/github.com/bytecodealliance/wasmtime-go@v0.33.0/build/macos-aarch64

wasmtime/target/release: wasmtime
	cd wasmtime && cargo build -p wasmtime-c-api --release

wasmtime:
	git clone https://github.com/bytecodealliance/wasmtime.git --depth 1
	cd wasmtime && git fetch --tags
	cd wasmtime && git checkout v0.33.0
	cd wasmtime && git submodule update --init



