package bootstrap

import (
	"fmt"
	"heroku-line-bot/src/util"
	errUtil "heroku-line-bot/src/util/error"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v2"
)

var (
	cfg *Config
)

func Get() (*Config, errUtil.IError) {
	if cfg == nil {
		errInfo := loadConfig()
		if errInfo != nil {
			return nil, errInfo
		}
	}
	return cfg, nil
}

func loadConfig() errUtil.IError {
	configName := GetEnvConfig()
	if configName == "" {
		configName = "master"
	}

	var errInfo errUtil.IError
	cfg, errInfo = loadYmlConfig(configName)
	if errInfo != nil {
		return errInfo
	}

	if errInfo := loadEnv(); errInfo != nil {
		return errInfo
	}

	return nil
}

func loadYmlConfig(fileName string) (*Config, errUtil.IError) {
	root, err := GetRootDirPath()
	if err != nil {
		return nil, errUtil.NewError(err)
	}
	configDir := filepath.Join("config")
	path := fmt.Sprintf("%s/%s/%s.yml", root, configDir, fileName)

	cfgBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errUtil.NewError(err)
	}

	cfg = &Config{}
	if err := yaml.Unmarshal(cfgBytes, cfg); err != nil {
		return nil, errUtil.NewError(err)
	}

	return cfg, nil
}

func loadEnv() errUtil.IError {
	if cfg == nil {
		cfg = &Config{}
	}

	if envStr := GetEnvPort(); envStr != "" {
		port, err := strconv.Atoi(envStr)
		if err != nil {
			return errUtil.NewError(err)
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
			return errUtil.NewError(err)
		}
	}

	if envStr := GetEnvRedisUrl(); envStr != "" {
		if err := cfg.ClubRedis.ScanUrl(envStr); err != nil {
			return errUtil.NewError(err)
		}
	}

	return nil
}

func GetRootDirPath() (string, error) {
	dir := GetEnvWorkDir()
	if dir == "" {
		dir = "HerokuLineBot"
	}
	return util.GetRootOf(dir)
}
