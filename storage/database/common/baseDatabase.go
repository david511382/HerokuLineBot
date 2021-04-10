package common

import (
	"time"

	"github.com/jinzhu/gorm"
)

type BaseDatabase struct {
	Read  *gorm.DB
	Write *gorm.DB
}

func (d *BaseDatabase) Begin() *gorm.DB {
	return d.Write.Begin()
}

func (d *BaseDatabase) SetConnection(maxIdleConns, maxOpenConns int, maxLifetime time.Duration) {
	if d.Read != nil {
		d.setConnection(d.Read, maxIdleConns, maxOpenConns, maxLifetime)
	}
	if d.Write != nil {
		d.setConnection(d.Write, maxIdleConns, maxOpenConns, maxLifetime)
	}
}

func (d *BaseDatabase) setConnection(db *gorm.DB, maxIdleConns, maxOpenConns int, maxLifetime time.Duration) {
	sqlDB := db.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(maxIdleConns)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(maxOpenConns)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(maxLifetime)
}

func (d *BaseDatabase) Dispose() error {
	if d.Read != nil {
		if err := d.Read.Close(); err != nil {
			return err
		}
	}

	if d.Write != nil {
		if err := d.Write.Close(); err != nil {
			return err
		}
	}

	return nil
}
