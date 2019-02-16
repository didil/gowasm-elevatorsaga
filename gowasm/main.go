package main

import (
	"syscall/js"
)

func main() {
	ch := make(chan struct{})

	js.Global().Get("console").Call("log", "Starting Wasm module ...")

	init := js.NewCallback(Init)
	defer init.Release()
	update := js.NewCallback(Update)
	defer update.Release()

	exitWasm := js.NewCallback(func(args []js.Value) {
		ch <- struct{}{}
	})
	defer exitWasm.Release()

	c := make(map[string]interface{})
	c["init"] = init
	c["update"] = update
	js.Global().Get("GoWasmBuilder").Set("codeObj", c)
	js.Global().Get("GoWasmBuilder").Set("exitWasm", exitWasm)

	<-ch
	js.Global().Get("console").Call("log", "Exiting Wasm module ...")
}
