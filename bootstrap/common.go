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

	fs *embed.FS
)

func Get() *Config {
	return cfg
}

func LoadFS(f *embed.FS) {
	fs = f
}

// ReadConfig read config from filepath
func LoadConfig(fileName string) (*Config, errLogic.IError) {
	f := fs
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

func LoadEnv() errLogic.IError {
	if cfg == nil {
		cfg = &Config{}
	}

	if envStr := os.Getenv("PORT"); envStr != "" {
		port, err := strconv.Atoi(envStr)
		if err != nil {
			return errLogic.NewError(err)
		}
		cfg.Server.Port = port
	}

	if envStr := os.Getenv("LINE_BOT_ADMIN_ID"); envStr != "" {
		cfg.LineBot.AdminID = envStr
	}
	if envStr := os.Getenv("LINE_BOT_ROOM_ID"); envStr != "" {
		cfg.LineBot.RoomID = envStr
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
			return errLogic.NewError(err)
		}
	}

	if envStr := os.Getenv("REDIS_URL"); envStr != "" {
		if err := cfg.ClubRedis.ScanUrl(envStr); err != nil {
			return errLogic.NewError(err)
		}
	}

	return nil
}
