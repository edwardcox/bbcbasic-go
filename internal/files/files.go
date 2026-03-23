package files

import (
	"os"
)

func SaveTextFile(filename string, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

func LoadTextFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
