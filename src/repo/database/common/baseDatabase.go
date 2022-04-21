package common

import (
	errUtil "heroku-line-bot/src/pkg/util/error"

	"gorm.io/gorm"
)

type SchemaCreator[Schema IBaseDatabase] func(connect Connect) Schema
type Connect func() (master, slave *gorm.DB, resultErr error)

type IBaseDatabase interface {
	IConnectionCreator
	Dispose() errUtil.IError
}

type BaseDatabase[Schema IBaseDatabase] struct {
	*MasterSlaveManager
	schemaCreator SchemaCreator[Schema]
}

func NewBaseDatabase[Schema IBaseDatabase](
	connectionCreator Connect,
	schemaCreator func(connect Connect) Schema,
) *BaseDatabase[Schema] {
	result := &BaseDatabase[Schema]{
		MasterSlaveManager: NewMasterSlaveManager(connectionCreator),
		schemaCreator:      schemaCreator,
	}
	return result
}

func (d *BaseDatabase[Schema]) Begin() (
	db Schema,
	trans ITransaction,
	resultErr error,
) {
	connect := func() (master *gorm.DB, slave *gorm.DB, resultErr error) {
		dp, err := d.GetMaster()
		if err != nil {
			resultErr = err
			return
		}

		dp = dp.Begin()
		if dp.Error != nil {
			resultErr = dp.Error
			return
		}

		master = dp
		slave = dp
		return
	}
	db = d.schemaCreator(connect)

	conn, err := db.GetMaster()
	if err != nil {
		resultErr = err
		return
	}

	trans = NewTransaction(conn)
	return
}
