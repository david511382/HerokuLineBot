package bootstrap

import (
	"os"
)

func GetEnvConfig() string {
	return os.Getenv("CONFIG")
}

func SetEnvConfig(s string) error {
	return os.Setenv("CONFIG", s)
}

func GetEnvPort() string {
	return os.Getenv("PORT")
}

func GetEnvLineBotAdminID() string {
	return os.Getenv("LINE_BOT_ADMIN_ID")
}

func GetEnvLineBotChannelAccessToken() string {
	return os.Getenv("LINE_BOT_CHANNEL_ACCESS_TOKEN")
}

func GetEnvTelegramBotAdminID() string {
	return os.Getenv("TELEGRAM_BOT_ADMIN_ID")
}

func GetEnvTelegramBotChannelAccessToken() string {
	return os.Getenv("TELEGRAM_BOT_CHANNEL_ACCESS_TOKEN")
}

func GetEnvGoogleScriptUrl() string {
	return os.Getenv("GOOGLE_SCRIPT_URL")
}

func GetEnvDatabaseUrl() string {
	return os.Getenv("DATABASE_URL")
}

func GetEnvRedisUrl() string {
	return os.Getenv("REDIS_URL")
}

func GetEnvWorkDir() string {
	return os.Getenv("WORK_DIR")
}

func SetEnvWorkDir(s string) error {
	return os.Setenv("WORK_DIR", s)
}
