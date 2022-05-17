package common

import (
	"heroku-line-bot/src/pkg/util"

	"gorm.io/gorm"
)

type SchemaCreator[Schema IBaseDatabase] func(connect func() (master *gorm.DB, slave *gorm.DB, resultErr error)) Schema

type IBaseDatabase interface {
	IConnection
	Dispose() error
}

type BaseDatabase[Schema IBaseDatabase] struct {
	*util.MasterSlaveManager[*gorm.DB]
	schemaCreator SchemaCreator[Schema]
}

func NewBaseDatabase[Schema IBaseDatabase](
	connectionCreator func() (master *gorm.DB, slave *gorm.DB, resultErr error),
	schemaCreator func(connect func() (master *gorm.DB, slave *gorm.DB, resultErr error)) Schema,
) *BaseDatabase[Schema] {
	result := &BaseDatabase[Schema]{
		MasterSlaveManager: util.NewMasterSlaveManager(
			connectionCreator,
			DisposeConnection,
		),
		schemaCreator: schemaCreator,
	}
	return result
}

func (d *BaseDatabase[Schema]) Begin() (
	db Schema,
	trans ITransaction,
	resultErr error,
) {
	connect := func() (master, slave *gorm.DB, resultErr error) {
		dp, err := d.MasterSlaveManager.GetMaster()
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
