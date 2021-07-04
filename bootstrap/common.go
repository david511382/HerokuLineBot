package bootstrap

import (
	"embed"
	errLogic "heroku-line-bot/logic/error"
	"heroku-line-bot/util"
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

// ReadConfig read config from filepath
func ReadConfig(f *embed.FS, fileName string) (*Config, *errLogic.ErrorInfo) {
	var cfgBytes []byte
	if f != nil {
		fileBs, err := f.ReadFile(fileName)
		if err != nil {
			return nil, errLogic.NewError(err)
		}
		cfgBytes = fileBs
	} else {
		fileBs, err := util.ReadFile(fileName)
		if err != nil {
			return nil, errLogic.NewError(err)
		}
		cfgBytes = fileBs
	}

	cfg = &Config{}
	if err := yaml.Unmarshal(cfgBytes, cfg); err != nil {
		return nil, errLogic.NewError(err)
	}

	return cfg, nil
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
