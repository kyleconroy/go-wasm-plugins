package main

import (
	"context"
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/trace"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
)

//go:embed plugin_rust.wasm
var wasmRust []byte

//go:embed plugin_go.wasm
var wasmGo []byte

//go:embed plugin_go_json.wasm
var wasmGoJSON []byte

// TODO: https://pkg.go.dev/github.com/bytecodealliance/wasmtime-go#example-package-Wasi

func run(cctx context.Context, name string, wasm []byte) error {
	ctx, task := trace.NewTask(cctx, name)
	defer task.End()

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

	moduRegion := trace.StartRegion(ctx, "wasmtime.NewModule")
	module, err := wasmtime.NewModule(store.Engine, wasm)
	moduRegion.End()
	if err != nil {
		return fmt.Errorf("define wasi: %w", err)
	}

	linkRegion := trace.StartRegion(ctx, "linker.Instantiate")
	instance, err := linker.Instantiate(store, module)
	linkRegion.End()
	if err != nil {
		return fmt.Errorf("define wasi: %w", err)
	}

	// Run the function

	callRegion := trace.StartRegion(ctx, "call _start")
	nom := instance.GetExport(store, "_start").Func()
	_, err = nom.Call(store)
	callRegion.End()
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
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace file: %v", err)
		}
	}()

	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()

	ctx := context.Background()

	fmt.Println(ctx, "rust")
	if e := run(ctx, "rust", wasmRust); e != nil {
		log.Fatal(e)
	}
	fmt.Println(ctx, "go-proto")
	if e := run(ctx, "go-proto", wasmGo); e != nil {
		log.Fatal(e)
	}
	fmt.Println(ctx, "go-json")
	if e := run(ctx, "go-json", wasmGoJSON); e != nil {
		log.Fatal(e)
	}
}
