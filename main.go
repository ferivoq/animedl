package main

import (
	"animedrive-dl/cmd"
	"embed"
	"encoding/json"
	"log"
)

//go:embed headers.json
var headersFile embed.FS

func loadHeaders() (map[string]map[string]string, error) {
	data, err := headersFile.ReadFile("headers.json")
	if err != nil {
		return nil, err
	}

	var headers map[string]map[string]string
	err = json.Unmarshal(data, &headers)
	if err != nil {
		return nil, err
	}

	return headers, nil
}

func main() {
	headers, err := loadHeaders()
	if err != nil {
		log.Fatalf("Failed to load headers: %v", err)
	}

	cmd.Execute(headers)
}
