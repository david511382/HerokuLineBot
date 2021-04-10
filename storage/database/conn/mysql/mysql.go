package mysql

import (
	"fmt"
	"heroku-line-bot/bootstrap"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type mysql struct {
	cfg bootstrap.Db
}

func (d mysql) Connect() (*gorm.DB, error) {
	addr := d.addr()
	return gorm.Open("mysql", addr)
}

func (d mysql) addr() string {
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
