package common

import (
	"heroku-line-bot/src/pkg/util"

	"github.com/go-redis/redis"
)

type BaseDatabase[Schema any] struct {
	*util.MasterSlaveManager[*redis.Client]
	schemaCreator func(connectionCreator IConnection) Schema
	baseKey       string
}

func NewBaseDatabase[Schema any](
	connectionCreator func() (master, slave *redis.Client, resultErr error),
	schemaCreator func(connectionCreator IConnection) Schema,
	baseKey string,
) *BaseDatabase[Schema] {
	result := &BaseDatabase[Schema]{
		MasterSlaveManager: util.NewMasterSlaveManager(
			connectionCreator,
			func(conn *redis.Client) error {
				return conn.Close()
			},
		),
		schemaCreator: schemaCreator,
		baseKey:       baseKey,
	}
	return result
}

func (d *BaseDatabase[Schema]) GetSlave() (redis.Cmdable, error) {
	return d.MasterSlaveManager.GetSlave()
}

func (d *BaseDatabase[Schema]) GetMaster() (redis.Cmdable, error) {
	return d.MasterSlaveManager.GetMaster()
}

func (d *BaseDatabase[Schema]) GetBaseKey() string {
	return d.baseKey
}

func (d *BaseDatabase[Schema]) Transaction() (
	db Schema,
	commitFn func() error,
) {
	commitFn = func() error { return nil }

	fn := func() (master, slave redis.Cmdable, resultErr error) {
		writeConn, err := d.MasterSlaveManager.GetMaster()
		if err != nil {
			resultErr = err
			return
		}
		readConn, err := d.MasterSlaveManager.GetSlave()
		if err != nil {
			resultErr = err
			return
		}

		pipe := writeConn.TxPipeline()
		master = pipe
		slave = readConn

		commitFn = func() error {
			if _, err := pipe.Exec(); err != nil {
				return err
			}

			return nil
		}

		return
	}
	db = d.schemaCreator(NewCmdableDatabase(fn))

	return
}

type CmdableDatabase struct {
	*util.MasterSlaveManager[redis.Cmdable]
}

func NewCmdableDatabase(
	connect func() (write, read redis.Cmdable, resultErr error),
) *CmdableDatabase {
	return &CmdableDatabase{
		MasterSlaveManager: util.NewMasterSlaveManager(
			connect,
			func(conn redis.Cmdable) error {
				return nil
			},
		),
	}
}
