package conn

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/repo/database/conn/mysql"
	"heroku-line-bot/src/repo/database/conn/postgre"
	"heroku-line-bot/src/repo/database/domain"
	"strings"

	"gorm.io/gorm"
)

func Connect(cfg bootstrap.Db) (*gorm.DB, error) {
	var c IConnect
	dbs := []domain.DbType{
		domain.POSTGRE_DB_TYPE,
		domain.MYSQL_DB_TYPE,
	}
	for _, protocol := range dbs {
		protocolStr := string(protocol)
		if strings.HasPrefix(cfg.Protocol, protocolStr) {
			switch protocol {
			case domain.POSTGRE_DB_TYPE:
				c = postgre.New(cfg)
			case domain.MYSQL_DB_TYPE:
				c = mysql.New(cfg)
			default:
				return nil, domain.UNKNOWN_DB_TYPE_ERROR
			}
			break
		}
	}
	if c == nil {
		return nil, domain.UNKNOWN_DB_TYPE_ERROR
	}

	dialector := c.GetDialector()
	return gorm.Open(dialector, &gorm.Config{})
}
