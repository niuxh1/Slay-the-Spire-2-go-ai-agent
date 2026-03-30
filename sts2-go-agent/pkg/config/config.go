package config

import (
	"encoding/json"
	"os"
)

type AppConfig struct {
	APIKey     string `json:"api_key"`
	BaseURL    string `json:"base_url"`
	Model      string `json:"model"`
	GameURL    string `json:"game_url"`
}

const configPath = "sts2_config.json"

func Load() AppConfig {
	var cfg AppConfig
	data, err := os.ReadFile(configPath)
	if err == nil {
		json.Unmarshal(data, &cfg)
	}
	// 默认值
	if cfg.GameURL == "" { cfg.GameURL = "http://127.0.0.1:8080" }
	if cfg.BaseURL == "" { cfg.BaseURL = "https://api.openai.com/v1" }
	if cfg.Model == "" { cfg.Model = "gpt-4o" }
	return cfg
}

func (c AppConfig) Save() {
	data, _ := json.MarshalIndent(c, "", "  ")
	os.WriteFile(configPath, data, 0644)
}
