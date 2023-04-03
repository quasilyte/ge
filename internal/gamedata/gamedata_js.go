package gamedata

import (
	"syscall/js"
)

func dataPath(appName, itemKey string) string {
	return "ge_game_" + appName + "_" + itemKey
}

func dataExists(appName, itemKey string) (bool, error) {
	result := js.Global().Get("localStorage").Call("getItem", dataPath(appName, itemKey))
	return !result.IsNull(), nil
}

func saveData(appName, itemKey string, data []byte) error {
	js.Global().Get("localStorage").Call("setItem", "ge_game_"+appName+"_"+itemKey, string(data))
	return nil
}

func loadData(appName, itemKey string) ([]byte, error) {
	result := js.Global().Get("localStorage").Call("getItem", "ge_game_"+appName+"_"+itemKey)
	if result.IsNull() {
		return nil, nil
	}
	return []byte(result.String()), nil
}
