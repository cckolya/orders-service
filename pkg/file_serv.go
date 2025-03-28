package pkg

import (
	"embed"
)

//go:embed data.txt
var fs embed.FS

func ReadFile() (string, error) {
	data, err := fs.ReadFile("data.txt")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
