package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	DefaultDB    string `json:"defaultDB"`
	FileSystemDB string `json:"fileSystemDB"`
}

func ReadFromFile(file string) (c Config, err error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, &c)
	return
}
