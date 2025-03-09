//go:build wasm
// +build wasm

package main

import (
	"syscall/js"

	"gotoproto/internal/converter"
)

func main() {
	js.Global().Set("Convert", js.FuncOf(func(this js.Value, args []js.Value) any {
		result := converter.Convert(args[0].String())
		doc := js.Global().Get("document")
		output := doc.Call("getElementById", "output")
		output.Set("innerHTML", result)

		return js.ValueOf(result)
	}))

	// block so the function will always be available
	ch := make(chan struct{})
	ch <- struct{}{}
}
