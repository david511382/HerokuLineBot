package bootstrap

import (
	"fmt"
	"heroku-line-bot/src/pkg/util"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v2"
)

const (
	DEFAULT_WORK_DIR  = "HerokuLineBot"
	DEFAULT_IANA_ZONE = "Asia/Taipei"
)

var (
	cfg *Config
)

func Get() (*Config, error) {
	if cfg == nil {
		err := loadConfig()
		if err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

func loadConfig() error {
	configName := GetEnvConfig()
	if configName == "" {
		configName = "master"
	}

	var err error
	cfg, err = loadYmlConfig(configName)
	if err != nil {
		return err
	}

	loadDefault()

	if err := loadEnv(); err != nil {
		return err
	}

	return nil
}

func loadYmlConfig(fileName string) (*Config, error) {
	root, err := GetRootDirPath()
	if err != nil {
		return nil, err
	}
	configDir := filepath.Join("config")
	path := fmt.Sprintf("%s/%s/%s.yml", root, configDir, fileName)

	cfgBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg = &Config{}
	if err := yaml.Unmarshal(cfgBytes, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadDefault() {
	if cfg == nil {
		cfg = &Config{}
	}

	if cfg.Var.TimeZone == "" {
		cfg.Var.TimeZone = DEFAULT_IANA_ZONE
	}
}

func loadEnv() error {
	if cfg == nil {
		cfg = &Config{}
	}

	if envStr := GetEnvPort(); envStr != "" {
		port, err := strconv.Atoi(envStr)
		if err != nil {
			return err
		}
		cfg.Server.Port = port
	}

	if envStr := GetEnvLineBotAdminID(); envStr != "" {
		cfg.LineBot.AdminID = envStr
	}
	if envStr := GetEnvLineBotChannelAccessToken(); envStr != "" {
		cfg.LineBot.ChannelAccessToken = envStr
	}

	if envStr := GetEnvTelegramBotAdminID(); envStr != "" {
		cfg.TelegramBot.AdminID = envStr
	}
	if envStr := GetEnvTelegramBotChannelAccessToken(); envStr != "" {
		cfg.TelegramBot.ChannelAccessToken = envStr
	}

	if envStr := GetEnvGoogleScriptUrl(); envStr != "" {
		cfg.GoogleScript.Url = envStr
	}

	if envStr := GetEnvDatabaseUrl(); envStr != "" {
		if err := cfg.ClubDb.ScanUrl(envStr); err != nil {
			return err
		}
	}

	if envStr := GetEnvRedisUrl(); envStr != "" {
		if err := cfg.ClubRedis.ScanUrl(envStr); err != nil {
			return err
		}
	}

	return nil
}

func GetRootDirPath() (string, error) {
	if dir := GetEnvWorkDir(); dir != "" {
		return util.GetRootOf(dir)
	}

	return ".", nil
}
