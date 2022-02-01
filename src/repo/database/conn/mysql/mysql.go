package mysql

import (
	"fmt"
	"heroku-line-bot/bootstrap"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlDb struct {
	cfg bootstrap.Db
}

func (d mysqlDb) GetDialector() gorm.Dialector {
	addr := d.addr()
	return mysql.Open(addr)
}

func (d mysqlDb) addr() string {
	cfg := d.cfg
	addr := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		cfg.User,
		cfg.Password,
		cfg.Addr(),
		cfg.Database,
		cfg.Param,
	)
	return addr
}
