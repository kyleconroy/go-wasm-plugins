.PHONY: run

run: cli/plugin.wasm cli/cli
	cd cli && ./cli

cli/plugin.wasm: plugin/src/main.rs plugin/Cargo.toml
	cd plugin && cargo build --target wasm32-wasi
	cp plugin/target/wasm32-wasi/debug/plugin.wasm cli/plugin.wasm

cli/cli: cli/main.go
	cd cli && CGO_CFLAGS="-I/Users/kyle/projects/go-wasm-plugins/wasmtime/crates/c-api/wasm-c-api/include -I/Users/kyle/projects/go-wasm-plugins/wasmtime/crates/c-api/include" CGO_LDFLAGS="-L/Users/kyle/projects/go-wasm-plugins/wasmtime/target/release/" go build ./...
