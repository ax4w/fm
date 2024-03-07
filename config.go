package main

import (
	"encoding/json"
	"os"
)

type configStruct struct {
	OpenInApp map[string]string `json:"open_in_app"`
	KeyBinds  map[string]string `json:"keybinds"`
}

func loadConfig() configStruct {
	f, err := os.ReadFile("./config.json")
	if err != nil {
		panic(err.Error())
	}
	var config configStruct
	json.Unmarshal(f, &config)
	return config
}
