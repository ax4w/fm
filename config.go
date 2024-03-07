package main

import (
	"encoding/json"
	"os"
)

type configStruct struct {
	OpenInApp map[string]string `json:"open_in_app"`
	KeyBinds  map[string]string `json:"keybinds"`
}

var config configStruct

func loadConfig() {
	f, _ := os.ReadFile("./config.json")
	json.Unmarshal(f, &config)
}
