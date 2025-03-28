package pkg

import "os"

func ReadFile() (string, error) {
	data, err := os.ReadFile("data.txt")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
