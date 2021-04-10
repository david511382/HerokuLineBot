package postgre

import (
	"fmt"
	"heroku-line-bot/bootstrap"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type postgre struct {
	cfg bootstrap.Db
}

func (d postgre) Connect() (*gorm.DB, error) {
	addr := d.addr()
	return gorm.Open("postgres", addr)
}

func (d postgre) addr() string {
	cfg := d.cfg
	addr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d %s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.Port,
		cfg.Param,
	)

	return addr
}
