package bootstrap

import (
	"embed"
	"fmt"
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
	fileName = fmt.Sprintf("config/%s.yml", fileName)
	cfgBytes, err := f.ReadFile(fileName)
	if err != nil {
		return err
	}

	cfg = &Config{}
	err = yaml.Unmarshal(cfgBytes, cfg)

	return err
}

func LoadEnv(cfg *Config) error {
	portStr := os.Getenv("PORT")
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return err
		}
		cfg.Server.Port = port
	}
	return nil
}
