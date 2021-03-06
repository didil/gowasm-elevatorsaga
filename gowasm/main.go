package main

import (
	"syscall/js"
)

func main() {
	// exit channel
	ch := make(chan struct{})

	js.Global().Get("console").Call("log", "Starting Wasm module ...")

	// init game callback
	init := js.NewCallback(Init)
	defer init.Release()
	// update game callback
	update := js.NewCallback(Update)
	defer update.Release()

	// exit callback
	exitWasm := js.NewCallback(func(args []js.Value) {
		ch <- struct{}{}
	})
	defer exitWasm.Release()

	// create js objects
	c := make(map[string]interface{})
	c["init"] = init
	c["update"] = update
	js.Global().Get("GoWasmBuilder").Set("codeObj", c)
	js.Global().Get("GoWasmBuilder").Set("exitWasm", exitWasm)

	// wait for exit signal
	<-ch
	js.Global().Get("console").Call("log", "Exiting Wasm module ...")
}
