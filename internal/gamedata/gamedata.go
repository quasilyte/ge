package gamedata

func Save(appName, itemKey string, data []byte) error {
	return saveData(appName, itemKey, data)
}

func Load(appName, itemKey string) ([]byte, error) {
	return loadData(appName, itemKey)
}

func Exists(appName, itemKey string) (bool, error) {
	return dataExists(appName, itemKey)
}

func Locate(appName, itemKey string) string {
	return dataPath(appName, itemKey)
}
