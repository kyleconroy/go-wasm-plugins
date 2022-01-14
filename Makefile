.PHONY: run

run: cli/plugin.wasm cli/cli
	cd cli && ./cli

cli/plugin.wasm: plugin/src/main.rs plugin/Cargo.toml
	cd plugin && cargo build --target wasm32-wasi
	cp plugin/target/wasm32-wasi/debug/plugin.wasm cli/plugin.wasm

cli/cli: cli/main.go wasmtime/target/release
	cd cli && CGO_CFLAGS="-I/Users/kyle/projects/go-wasm-plugins/wasmtime/crates/c-api/wasm-c-api/include -I/Users/kyle/projects/go-wasm-plugins/wasmtime/crates/c-api/include" CGO_LDFLAGS="-L/Users/kyle/projects/go-wasm-plugins/wasmtime/target/release/" go build ./...

wasmtime/target/release: wasmtime
	cd wasmtime && cargo build -p wasmtime-c-api --release

wasmtime:
	git clone https://github.com/bytecodealliance/wasmtime.git --depth 1
	cd wasmtime && git fetch --tags
	cd wasmtime && git checkout v0.33.0
	cd wasmtime && git submodule update --init



