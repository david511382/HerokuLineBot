package common

import (
	errUtil "heroku-line-bot/src/pkg/util/error"
	"sync"

	"gorm.io/gorm"
)

type MasterSlaveManager struct {
	read              *gorm.DB
	write             *gorm.DB
	connectionCreator Connect
	sync.RWMutex
}

func NewMasterSlaveManager(
	connectionCreator Connect,
) *MasterSlaveManager {
	result := &MasterSlaveManager{
		connectionCreator: connectionCreator,
	}
	return result
}

func (d *MasterSlaveManager) connect() error {
	d.Lock()
	defer d.Unlock()
	if d.read != nil && d.write != nil {
		return nil
	}

	master, slave, err := d.connectionCreator()
	if err != nil {
		return err
	}
	if d.read == nil {
		d.read = slave
	} else {
		_ = disposeConnection(slave)
	}
	if d.write == nil {
		d.write = master
	} else {
		_ = disposeConnection(master)
	}
	return nil
}

func (d *MasterSlaveManager) GetSlave() (*gorm.DB, error) {
	d.RLock()
	isNoConnection := d.read == nil
	d.RUnlock()

	if isNoConnection {
		if err := d.connect(); err != nil {
			return nil, err
		}
	}
	return d.read, nil
}

func (d *MasterSlaveManager) GetMaster() (*gorm.DB, error) {
	d.RLock()
	isNoConnection := d.write == nil
	d.RUnlock()

	if isNoConnection {
		if err := d.connect(); err != nil {
			return nil, err
		}
	}
	return d.write, nil
}

func (d *MasterSlaveManager) Dispose() errUtil.IError {
	d.Lock()
	defer d.Unlock()

	if conn := d.read; conn != nil {
		if err := disposeConnection(conn); err != nil {
			return errUtil.NewError(err)
		}
	}

	if conn := d.write; conn != nil {
		if err := disposeConnection(conn); err != nil {
			return errUtil.NewError(err)
		}
	}

	return nil
}

func disposeConnection(conn *gorm.DB) error {
	sqlDB, err := conn.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}
	return nil
}
