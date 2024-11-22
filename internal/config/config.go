package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DBUrl string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}


func getConfigFilePath() (string, error) {
	fileLocation, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	fileLocation += "/.gatorconfig.json"

	return fileLocation, nil
}


func (c *Config) SetUser(username string) error {
	c.CurrentUsername = username
	fileLocation, err := getConfigFilePath()
	if err != nil {
		 return err
	}
	
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(fileLocation, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func Read() (Config, error) {

	var config Config;

	fileLocation, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(fileLocation)
	if err != nil {
		return Config{}, err
	}
	
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err	
	}
	return config, nil
}
