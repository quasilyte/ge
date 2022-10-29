package gamedata

import (
	"errors"
	"os"
	"path/filepath"
)

func getWindowsItemFolder(appName, itemKey string) (string, error) {
	appData := os.Getenv("AppData")
	if appData == "" {
		return "", errors.New("AppData env var is undefined")
	}
	dataPath := filepath.Join(appData, "ge_game_"+appName)
	if err := mkdirAll(dataPath); err != nil {
		return "", err
	}
	return filepath.Join(dataPath, itemKey+".json"), nil
}

func dataExists(appName, itemKey string) (bool, error) {
	itemPath, err := getWindowsItemFolder(appName, itemKey)
	if err != nil {
		return false, err
	}
	return fileExists(itemPath), nil
}

func saveData(appName, itemKey string, data []byte) error {
	itemPath, err := getWindowsItemFolder(appName, itemKey)
	if err != nil {
		return err
	}
	return os.WriteFile(itemPath, data, 0o666)
}

func loadData(appName, itemKey string) ([]byte, error) {
	itemPath, err := getWindowsItemFolder(appName, itemKey)
	if err != nil {
		return nil, err
	}
	if !fileExists(itemPath) {
		return nil, nil
	}
	return os.ReadFile(itemPath)
}
