//go:build wasm

package userdevice

import (
	"syscall/js"
)

func GetInfo() Info {
	var result Info
	result.IsMobile = js.Global().Get("matchMedia").Call("(hover: none)").Get("matches").Bool()
	return result
}
