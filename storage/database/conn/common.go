package conn

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/storage/database/conn/mysql"
	"heroku-line-bot/storage/database/conn/postgre"
	"heroku-line-bot/storage/database/domain"

	"github.com/jinzhu/gorm"
)

func Connect(cfg bootstrap.Db) (*gorm.DB, error) {
	var c IConnect
	dbType := cfg.Type
	switch dbType {
	case domain.POSTGRE_DB_TYPE:
		c = postgre.New(cfg)
	case domain.MYSQL_DB_TYPE:
		c = mysql.New(cfg)
	default:
		return nil, domain.UNKNOWN_DB_TYPE_ERROR
	}

	return c.Connect()
}
