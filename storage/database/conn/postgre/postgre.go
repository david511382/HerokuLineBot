package postgre

import (
	"fmt"
	"heroku-line-bot/bootstrap"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgreDb struct {
	cfg bootstrap.Db
}

func (d postgreDb) GetDialector() gorm.Dialector {
	addr := d.addr()
	return postgres.Open(addr)
}

func (d postgreDb) addr() string {
	cfg := d.cfg
	addr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s %s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.Param,
	)
	if cfg.Port > 0 {
		addr = fmt.Sprintf("%s port=%d",
			addr,
			cfg.Port,
		)
	}

	return addr
}
