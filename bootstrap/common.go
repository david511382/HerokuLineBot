package bootstrap

import (
	"embed"
	"fmt"
	errLogic "heroku-line-bot/logic/error"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

var (
	cfg *Config
)

func Get() *Config {
	return cfg
}

func LoadConfig(f embed.FS, fileName string) *Config {
	err := ReadConfig(f, fileName)
	if err != nil {
		panic(err)
	}
	return cfg
}

// ReadConfig read config from filepath
func ReadConfig(f embed.FS, fileName string) error {
	fileName = fmt.Sprintf("resource/config/%s.yml", fileName)
	cfgBytes, err := f.ReadFile(fileName)
	if err != nil {
		return err
	}

	cfg = &Config{}
	err = yaml.Unmarshal(cfgBytes, cfg)

	return err
}

func LoadEnv(cfg *Config) *errLogic.ErrorInfo {
	portStr := os.Getenv("PORT")
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return errLogic.NewError(err)
		}
		cfg.Server.Port = port
	}
	return nil
}
