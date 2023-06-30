package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// main is an example of how to extend a Go application with an addition
// function defined in WebAssembly.

//go:embed python-3.11.1.wasm
var python []byte

func main() {

	for {
		runme()
		runtime.GC()
		time.Sleep(time.Second)
	}
}

func runme() {
	ctx, cncl := context.WithCancel(context.Background())
	rConfig := wazero.NewRuntimeConfig().
		WithMemoryLimitPages(32768).
		WithCloseOnContextDone(true)
	r := wazero.NewRuntimeWithConfig(ctx, rConfig)

	apiCloser, err := wasi_snapshot_preview1.Instantiate(ctx, r)
	mod, err := r.CompileModule(ctx, python)
	if err != nil {
		log.Panicf("failed to instantiate module: %v", err)
	}
	fmt.Println(mod)
	cncl()
	// Note, pretty sure this will always return immediately if cancel has
	// returned since the cancellation sets the done flag inside the context
	// before returning.
	<-ctx.Done()
	err = apiCloser.Close(ctx)
	if err != nil {
		log.Panicf("failed to close module: %v", err)
	}
	err = r.Close(ctx)
	if err != nil {
		log.Panicf("failed to close runtime: %v", err)
	}
}
