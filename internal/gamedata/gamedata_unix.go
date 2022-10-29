//go:build darwin || linux

package gamedata

import (
	"os"
	"path/filepath"
)

func getUnixItemFolder(appName, itemKey string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dataPath := filepath.Join(home, ".local", "share", "ge_game_"+appName)
	if err := mkdirAll(dataPath); err != nil {
		return "", err
	}
	return filepath.Join(dataPath, itemKey+".json"), nil
}

func dataExists(appName, itemKey string) (bool, error) {
	itemPath, err := getUnixItemFolder(appName, itemKey)
	if err != nil {
		return false, err
	}
	return fileExists(itemPath), nil
}

func saveData(appName, itemKey string, data []byte) error {
	itemPath, err := getUnixItemFolder(appName, itemKey)
	if err != nil {
		return err
	}
	return os.WriteFile(itemPath, data, 0o666)
}

func loadData(appName, itemKey string) ([]byte, error) {
	itemPath, err := getUnixItemFolder(appName, itemKey)
	if err != nil {
		return nil, err
	}
	if !fileExists(itemPath) {
		return nil, nil
	}
	return os.ReadFile(itemPath)
}
