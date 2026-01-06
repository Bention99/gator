package config

import (
	"os"
	"encoding/json"
)

const configFileName = ".gatorconfig.json"

func Read() Config {
	filePath := os.UserHomeDir
	data, err := os.ReadFile(filePath+configFileName)
	if err != nil {
		fmt.Println("Error reading the json File. Error:")
		fmt.Printf(" - %v\n", err)
		return Config{}
	}
	c, err := unmarshal[Config](data)
	if err != nil {
		fmt.Println("Error unmarshaling the json File. Error:")
		fmt.Printf(" - %v\n", err)
		return Config{}
	}
	return c
}

func unmarshal[T any](b []byte) (T, error) {
	var v T
	if err := json.Unmarshal(b, &v); err != nil {
		var zero T
		return zero, err
	}
	return v, nil
}