package main

import (
	"strings"
	"syscall/js"
	"text/template"

	"github.com/ddddddO/ps2"
)

func main() {
	c := make(chan struct{}, 0)
	println("ps2 WebAssembly Initialized")
	registerCallbacks()
	<-c
}

func registerCallbacks() {
	js.Global().Set("ps2Run", js.FuncOf(ps2Run))
}

func ps2Run(this js.Value, args []js.Value) interface{} {
	document := js.Global().Get("document")
	getElementByID := getElementByIDFunc(document)

	input := getElementByID("input").Get("value").String()

	outputOfOption := ps2.WithOutputTypeJSON()
	outputRadios := getElementsByNameFunc(document)("output_type")
	for i := 0; i < outputRadios.Length(); i++ {
		radio := outputRadios.Index(i)
		if radio.Get("checked").Bool() {
			switch radio.Get("value").String() {
			case "json":
				outputOfOption = ps2.WithOutputTypeJSON()
			case "yaml":
				outputOfOption = ps2.WithOutputTypeYAML()
			case "toml":
				outputOfOption = ps2.WithOutputTypeTOML()
			default:
				outputOfOption = ps2.WithOutputTypeJSON()
			}
		}
	}

	output, err := ps2.Run(strings.NewReader(input), outputOfOption)
	if err != nil {
		alert(err.Error())
		return nil
	}

	div := getElementByID("result")
	if prePre := getElementByID("redered_json"); !prePre.IsNull() {
		removeChildFunc(div)(prePre)
	}

	pre := createElementFunc(document)("pre")
	pre.Set("id", "redered_json")
	pre.Set("innerHTML", template.HTMLEscapeString(output))
	appendChildFunc(div)(pre)
	appendChildFunc(getElementByID("main"))(div)

	return nil
}

func getElementByIDFunc(document js.Value) func(id string) js.Value {
	return func(id string) js.Value {
		return document.Call("getElementById", id)
	}
}

func getElementsByNameFunc(document js.Value) func(id string) js.Value {
	return func(id string) js.Value {
		return document.Call("getElementsByName", id)
	}
}

func createElementFunc(document js.Value) func(element string) js.Value {
	return func(element string) js.Value {
		return document.Call("createElement", element)
	}
}

func removeChildFunc(element js.Value) func(target js.Value) {
	return func(target js.Value) {
		element.Call("removeChild", target)
	}
}

func appendChildFunc(element js.Value) func(target js.Value) {
	return func(target js.Value) {
		element.Call("appendChild", target)
	}
}

func alert(msg string) {
	js.Global().Call("alert", msg)
}
