package common

import (
	errUtil "heroku-line-bot/src/pkg/util/error"
	"time"

	"gorm.io/gorm"
)

type BaseDatabase struct {
	read  *gorm.DB
	write *gorm.DB
}

func NewBaseDatabase(read, write *gorm.DB) *BaseDatabase {
	result := &BaseDatabase{
		read:  read,
		write: write,
	}
	return result
}

func (d *BaseDatabase) GetSlave() *gorm.DB {
	return d.read
}

func (d *BaseDatabase) GetMaster() *gorm.DB {
	return d.write
}

func (d *BaseDatabase) SetConnection(maxIdleConns, maxOpenConns int, maxLifetime time.Duration) errUtil.IError {
	if d.read != nil {
		if errInfo := d.setConnection(d.read, maxIdleConns, maxOpenConns, maxLifetime); errInfo != nil {
			return errInfo
		}
	}
	if d.write != nil {
		if errInfo := d.setConnection(d.write, maxIdleConns, maxOpenConns, maxLifetime); errInfo != nil {
			return errInfo
		}
	}

	return nil
}

func (d *BaseDatabase) setConnection(db *gorm.DB, maxIdleConns, maxOpenConns int, maxLifetime time.Duration) errUtil.IError {
	sqlDB, err := db.DB()
	if err != nil {
		return errUtil.NewError(err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(maxIdleConns)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(maxOpenConns)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(maxLifetime)

	return nil
}

func (d *BaseDatabase) Dispose() errUtil.IError {
	if d.read != nil {
		sqlDB, err := d.read.DB()
		if err != nil {
			return errUtil.NewError(err)
		}

		if err := sqlDB.Close(); err != nil {
			return errUtil.NewError(err)
		}
	}

	if d.write != nil {
		sqlDB, err := d.write.DB()
		if err != nil {
			return errUtil.NewError(err)
		}

		if err := sqlDB.Close(); err != nil {
			return errUtil.NewError(err)
		}
	}

	return nil
}
