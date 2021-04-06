package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
)

// TODO: https://pkg.go.dev/github.com/bytecodealliance/wasmtime-go#example-package-Wasi

func run() error {
	dir, err := ioutil.TempDir("", "out")
	if err != nil {
		return fmt.Errorf("temp dir: %w", err)
	}

	defer os.RemoveAll(dir)
	stdinPath := filepath.Join(dir, "stdin")
	stderrPath := filepath.Join(dir, "stderr")
	stdoutPath := filepath.Join(dir, "stdout")

	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)
	linker := wasmtime.NewLinker(store)

	// Configure WASI imports to write stdout into a file.
	wasiConfig := wasmtime.NewWasiConfig()
	wasiConfig.SetStdinFile(stdinPath)
	wasiConfig.SetStdoutFile(stdoutPath)
	wasiConfig.SetStderrFile(stderrPath)

	// Set the version to the same as in the WAT.
	wasi, err := wasmtime.NewWasiInstance(store, wasiConfig, "wasi_snapshot_preview1")
	if err != nil {
		return fmt.Errorf("new wasi instances: %w", err)
	}

	// Link WASI
	err = linker.DefineWasi(wasi)
	if err != nil {
		return fmt.Errorf("define wasi: %w", err)
	}

	// Create our module
	//
	// Compiling modules requires WebAssembly binary input, but the wasmtime
	// package also supports converting the WebAssembly text format to the
	// binary format.
	wasm, err := os.ReadFile("plugin.wasm")
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	module, err := wasmtime.NewModule(store.Engine, wasm)
	if err != nil {
		return fmt.Errorf("define wasi: %w", err)
	}

	instance, err := linker.Instantiate(module)
	if err != nil {
		return fmt.Errorf("define wasi: %w", err)
	}

	// Run the function
	nom := instance.GetExport("_start").Func()
	_, err = nom.Call()
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
	if e := run(); e != nil {
		log.Fatal(e)
	}
}
