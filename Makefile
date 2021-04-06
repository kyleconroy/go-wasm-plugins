.PHONY: run

run: cli/plugin.wasm cli/cli
	cd cli && ./cli

cli/plugin.wasm:
	cd plugin && cargo build --target wasm32-wasi
	cp plugin/target/wasm32-wasi/debug/plugin.wasm cli/plugin.wasm

cli/cli: cli/main.go
	cd cli && go build ./...


