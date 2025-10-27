package file

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func SaveJSON(filePath string, data interface{}) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0644)
}
