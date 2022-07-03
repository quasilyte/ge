//go:build wasm

package locales

import (
	"syscall/js"
)

func inferLanguages() []string {
	jsNavigator := js.Global().Get("navigator")

	var langs []string
	lang := jsNavigator.Get("language").String()
	if lang != "" {
		langs = append(langs, lang)
	}
	return langs
}
