package common

import (
	"time"

	errLogic "heroku-line-bot/logic/error"

	"gorm.io/gorm"
)

type BaseDatabase struct {
	Read  *gorm.DB
	Write *gorm.DB
}

func (d *BaseDatabase) Begin() *gorm.DB {
	return d.Write.Begin()
}

func (d *BaseDatabase) SetConnection(maxIdleConns, maxOpenConns int, maxLifetime time.Duration) errLogic.IError {
	if d.Read != nil {
		if errInfo := d.setConnection(d.Read, maxIdleConns, maxOpenConns, maxLifetime); errInfo != nil {
			return errInfo
		}
	}
	if d.Write != nil {
		if errInfo := d.setConnection(d.Read, maxIdleConns, maxOpenConns, maxLifetime); errInfo != nil {
			return errInfo
		}
	}

	return nil
}

func (d *BaseDatabase) setConnection(db *gorm.DB, maxIdleConns, maxOpenConns int, maxLifetime time.Duration) errLogic.IError {
	sqlDB, err := db.DB()
	if err != nil {
		return errLogic.NewError(err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(maxIdleConns)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(maxOpenConns)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(maxLifetime)

	return nil
}

func (d *BaseDatabase) Dispose() errLogic.IError {
	if d.Read != nil {
		sqlDB, err := d.Read.DB()
		if err != nil {
			return errLogic.NewError(err)
		}

		if err := sqlDB.Close(); err != nil {
			return errLogic.NewError(err)
		}
	}

	if d.Write != nil {
		sqlDB, err := d.Write.DB()
		if err != nil {
			return errLogic.NewError(err)
		}

		if err := sqlDB.Close(); err != nil {
			return errLogic.NewError(err)
		}
	}

	return nil
}
