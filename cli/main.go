package main

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
)

//go:embed plugin_rust.wasm
var wasmRust []byte

//go:embed plugin_go.wasm
var wasmGo []byte

// TODO: https://pkg.go.dev/github.com/bytecodealliance/wasmtime-go#example-package-Wasi

func run(wasm []byte) error {
	dir, err := ioutil.TempDir("", "out")
	if err != nil {
		return fmt.Errorf("temp dir: %w", err)
	}

	defer os.RemoveAll(dir)
	stdinPath := filepath.Join(dir, "stdin")
	stderrPath := filepath.Join(dir, "stderr")
	stdoutPath := filepath.Join(dir, "stdout")

	engine := wasmtime.NewEngine()
	linker := wasmtime.NewLinker(engine)

	// Link WASI
	if err := linker.DefineWasi(); err != nil {
		return fmt.Errorf("define wasi: %w", err)
	}

	// Configure WASI imports to write stdout into a file.
	wasiConfig := wasmtime.NewWasiConfig()
	wasiConfig.SetStdinFile(stdinPath)
	wasiConfig.SetStdoutFile(stdoutPath)
	wasiConfig.SetStderrFile(stderrPath)

	store := wasmtime.NewStore(engine)
	store.SetWasi(wasiConfig)

	// Set the version to the same as in the WAT.
	// wasi, err := wasmtime.NewWasiInstance(store, wasiConfig, "wasi_snapshot_preview1")
	// if err != nil {
	// 	return fmt.Errorf("new wasi instances: %w", err)
	// }

	// Create our module
	//
	// Compiling modules requires WebAssembly binary input, but the wasmtime
	// package also supports converting the WebAssembly text format to the
	// binary format.
	// wasm, err := os.ReadFile("plugin.wasm")
	// if err != nil {
	// 	return fmt.Errorf("read file: %w", err)
	// }

	module, err := wasmtime.NewModule(store.Engine, wasm)
	if err != nil {
		return fmt.Errorf("define wasi: %w", err)
	}

	instance, err := linker.Instantiate(store, module)
	if err != nil {
		return fmt.Errorf("define wasi: %w", err)
	}

	// Run the function
	nom := instance.GetExport(store, "_start").Func()
	_, err = nom.Call(store)
	if err != nil {
		return fmt.Errorf("call: %w", err)
	}

	// Print WASM stdout
	out, err := os.ReadFile(stdoutPath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	fmt.Print(string(out))
	return nil
}

func main() {
	fmt.Println("rust")
	if e := run(wasmRust); e != nil {
		log.Fatal(e)
	}
	fmt.Println("go")
	if e := run(wasmGo); e != nil {
		log.Fatal(e)
	}

}
