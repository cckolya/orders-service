package main

import (
	"embed"
	"fmt"
)

//go:embed data.txt
var fs embed.FS

func main() {
	fmt.Println(ReadFile())
}

func ReadFile() (string, error) {
	data, err := fs.ReadFile("data.txt")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
