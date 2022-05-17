package util

import (
	"sync"
)

type MasterSlaveManager[connection any] struct {
	read              connection
	write             connection
	connectionCreator func() (write, read connection, resultErr error)
	connectionCloser  func(conn connection) error
	sync.RWMutex
}

func NewMasterSlaveManager[connection any](
	connectionCreator func() (write, read connection, resultErr error),
	connectionCloser func(conn connection) error,
) *MasterSlaveManager[connection] {
	result := &MasterSlaveManager[connection]{
		connectionCreator: connectionCreator,
		connectionCloser:  connectionCloser,
	}
	return result
}

func (d *MasterSlaveManager[connection]) connect() error {
	d.Lock()
	defer d.Unlock()
	if !IsZero(d.read) &&
		!IsZero(d.write) {
		return nil
	}

	write, read, err := d.connectionCreator()
	if err != nil {
		return err
	}
	if IsZero(d.read) {
		d.read = read
	} else {
		_ = d.connectionCloser(read)
	}
	if IsZero(d.write) {
		d.write = write
	} else {
		_ = d.connectionCloser(write)
	}
	return nil
}

func (d *MasterSlaveManager[connection]) GetSlave() (connection, error) {
	d.RLock()
	isNoConnection := IsZero(d.read)
	d.RUnlock()

	if isNoConnection {
		if isNoConnection {
			if err := d.connect(); err != nil {
				return ZeroOf[connection](), err
			}
		}
	}
	return d.read, nil
}

func (d *MasterSlaveManager[connection]) GetMaster() (connection, error) {
	d.RLock()
	isNoConnection := IsZero(d.write)
	d.RUnlock()

	if isNoConnection {
		if isNoConnection {
			if err := d.connect(); err != nil {
				return ZeroOf[connection](), err
			}
		}
	}
	return d.write, nil
}

func (d *MasterSlaveManager[connection]) Dispose() error {
	d.Lock()
	defer d.Unlock()

	if conn := d.read; !IsZero(conn) {
		if err := d.connectionCloser(conn); err != nil {
			return err
		}
	}

	if conn := d.write; !IsZero(conn) {
		if err := d.connectionCloser(conn); err != nil {
			return err
		}
	}

	return nil
}
