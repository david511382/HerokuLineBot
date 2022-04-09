package common

import (
	errUtil "heroku-line-bot/src/pkg/util/error"

	"gorm.io/gorm"
)

type Connect func() (master, slave *gorm.DB, resultErr error)

type IBaseDatabase interface {
	GetSlave() *gorm.DB
	GetMaster() *gorm.DB
	Dispose() errUtil.IError
	BeginTransaction(
		creator func(connect Connect) (IBaseDatabase, error),
	) (
		trans ITransaction,
		resultErr error,
	)
}

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

func (d *BaseDatabase) Dispose() errUtil.IError {
	if d == nil {
		return nil
	}

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

func (d *BaseDatabase) BeginTransaction(
	creator func(connect Connect) (IBaseDatabase, error),
) (
	trans ITransaction,
	resultErr error,
) {
	connect := func() (master *gorm.DB, slave *gorm.DB, resultErr error) {
		dp := d.GetMaster().Begin()
		if dp.Error != nil {
			resultErr = dp.Error
			return
		}

		master = dp
		slave = dp
		return
	}
	db, err := creator(connect)
	if err != nil {
		resultErr = err
		return
	}

	trans = NewTransaction(db.GetMaster())
	return
}
