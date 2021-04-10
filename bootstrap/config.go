package bootstrap

import (
	databaseDomain "heroku-line-bot/storage/database/domain"
	"strconv"
)

type Config struct {
	Server       Server       `yaml:"server"`
	LineBot      LineBot      `yaml:"line_bot"`
	GoogleScript GoogleScript `yaml:"google_script"`
	DbConfig     DbConfig     `yaml:"db"`
	ClubDb       Db           `yaml:"club_db"`
	RedisConfig  DbConfig     `yaml:"redis"`
	ClubRedis    Db           `yaml:"club_redis"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (c *Server) Addr() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

type LineBot struct {
	AdminID            string `yaml:"admin_id"`
	RoomID             string `yaml:"room_id"`
	ChannelAccessToken string `yaml:"channel_access_token"`
}

type GoogleScript struct {
	Url string `yaml:"url"`
}

type DbConfig struct {
	MaxIdleConns int `yaml:"max_idle_conns"`
	MaxOpenConns int `yaml:"max_open_conns"`
	MaxLifeHour  int `yaml:"max_lifehour"`
}

type Db struct {
	Server   `yaml:"server"`
	Password string                `yaml:"password"`
	Database string                `yaml:"database"`
	User     string                `yaml:"user"`
	Type     databaseDomain.DbType `yaml:"type"`
	Param    string                `yaml:"param"`
}
