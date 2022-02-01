package common

import (
	"time"

	"github.com/go-redis/redis"
)

type IDatabase interface {
	InitModel(read, write redis.Cmdable)
}

type BaseDatabase struct {
	read    *redis.Client
	write   *redis.Client
	db      IDatabase
	baseKey string
}

func NewBaseDatabase(read, write *redis.Client, db IDatabase, baseKey string) *BaseDatabase {
	result := &BaseDatabase{
		read:    read,
		write:   write,
		db:      db,
		baseKey: baseKey,
	}
	return result
}

func (db *BaseDatabase) Init() {
	db.db.InitModel(db.read, db.write)
}

func (db *BaseDatabase) GetBaseKey() string {
	return db.baseKey
}

func (db *BaseDatabase) SetConnection(maxConnAge time.Duration) {
	if db.read != nil {
		db.setConnection(db.read, maxConnAge)
	}
	if db.write != nil {
		db.setConnection(db.write, maxConnAge)
	}
}

func (db *BaseDatabase) setConnection(connection *redis.Client, maxConnAge time.Duration) {
	connection.Options().MaxConnAge = maxConnAge
}

func (db *BaseDatabase) Dispose() error {
	if db.read != nil {
		if err := db.read.Close(); err != nil {
			return err
		}
	}

	if db.write != nil {
		if err := db.write.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (db *BaseDatabase) Transaction() (
	redisDb *BaseDatabase,
	commitF, rollbackF func() error,
) {
	pipe := db.write.TxPipeline()

	redisDb = NewBaseDatabase(db.read, db.write, db.db, db.baseKey)
	redisDb.db.InitModel(db.read, pipe)

	commitF = func() error {
		if _, err := pipe.Exec(); err != nil {
			return err
		}

		return nil
	}
	rollbackF = func() error {
		return pipe.Close()
	}

	return
}
