package file

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

var jsonWriteMu sync.Mutex

func SaveJSON(filePath string, data interface{}) error {
	jsonWriteMu.Lock()
	defer jsonWriteMu.Unlock()

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	tempFile := filePath + ".tmp"
	if err := os.WriteFile(tempFile, jsonData, 0644); err != nil {
		return err
	}

	return os.Rename(tempFile, filePath) // atomic only linux!!
}
