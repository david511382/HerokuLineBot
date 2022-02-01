package bootstrap

import (
	"fmt"
	"heroku-line-bot/util"
	errUtil "heroku-line-bot/util/error"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v2"
)

var (
	cfg *Config
)

func Get() *Config {
	return cfg
}

func SetEnvConfig(s string) error {
	return os.Setenv("CONFIG", s)
}

// ReadConfig read config from filepath
func LoadConfig() (*Config, errUtil.IError) {
	configName := os.Getenv("CONFIG")
	if configName == "" {
		configName = "master"
	}

	cfg = &Config{}

	return loadConfig(configName)
}

func loadConfig(fileName string) (*Config, errUtil.IError) {
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

func LoadEnv() errUtil.IError {
	if cfg == nil {
		cfg = &Config{}
	}

	if envStr := os.Getenv("PORT"); envStr != "" {
		port, err := strconv.Atoi(envStr)
		if err != nil {
			return errUtil.NewError(err)
		}
		cfg.Server.Port = port
	}

	if envStr := os.Getenv("LINE_BOT_ADMIN_ID"); envStr != "" {
		cfg.LineBot.AdminID = envStr
	}
	if envStr := os.Getenv("LINE_BOT_CHANNEL_ACCESS_TOKEN"); envStr != "" {
		cfg.LineBot.ChannelAccessToken = envStr
	}

	if envStr := os.Getenv("TELEGRAM_BOT_ADMIN_ID"); envStr != "" {
		cfg.TelegramBot.AdminID = envStr
	}
	if envStr := os.Getenv("TELEGRAM_BOT_CHANNEL_ACCESS_TOKEN"); envStr != "" {
		cfg.TelegramBot.ChannelAccessToken = envStr
	}

	if envStr := os.Getenv("GOOGLE_SCRIPT_URL"); envStr != "" {
		cfg.GoogleScript.Url = envStr
	}

	if envStr := os.Getenv("DATABASE_URL"); envStr != "" {
		if err := cfg.ClubDb.ScanUrl(envStr); err != nil {
			return errUtil.NewError(err)
		}
	}

	if envStr := os.Getenv("REDIS_URL"); envStr != "" {
		if err := cfg.ClubRedis.ScanUrl(envStr); err != nil {
			return errUtil.NewError(err)
		}
	}

	return nil
}

func GetRootDirPath() (string, error) {
	dir := cfg.Var.WorkDir
	if dir == "" {
		dir = "HerokuLineBot"
	}
	return util.GetRootOf(dir)
}
