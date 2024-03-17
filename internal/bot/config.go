package bot

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Token  string `json:"token"`
	ChatID int    `json:"chat_id"`
}

func ReadConfig() (Config, error) {
	bytes, err := os.ReadFile(getConfigPath())
	if err != nil {
		return Config{}, err
	}

	var conf Config

	if err = json.Unmarshal(bytes, &conf); err != nil {
		return Config{}, err
	}

	return conf, nil
}

func SaveConfig(conf Config) error {
	bytes, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	confPath := getConfigPath()

	if err := os.MkdirAll(filepath.Dir(confPath), 0770); err != nil {
		return err
	}

	return os.WriteFile(confPath, bytes, 0600)
}

func getConfigPath() string {
	configPath := os.Getenv("HOME")
	if configPath != "" {
		configPath = filepath.Join(configPath, ".config", "tg-agent-bot", "conf")
	} else {
		configPath = filepath.Join("/", "etc", "tg-agent-bot", "conf")
	}

	return configPath
}
