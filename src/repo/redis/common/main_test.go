package common

import (
	"fmt"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/pkg/test"
	"heroku-line-bot/src/repo/redis/conn"
	"os"
	"testing"

	"github.com/go-redis/redis"
)

var (
	defaultParser = newTestParser(
		func(key string) string {
			return key
		},
		func(keyStr string) (string, error) {
			return keyStr, nil
		},
		func(field string) string {
			return field
		},
		func(fieldStr string) (string, error) {
			return fieldStr, nil
		},
		func(value string) (string, error) {
			return value, nil
		},
		func(valueStr string) (string, error) {
			return valueStr, nil
		},
	)
)

func TestMain(m *testing.M) {
	if err := bootstrap.SetEnvWorkDir(bootstrap.DEFAULT_WORK_DIR); err != nil {
		panic(err)
	}
	if err := bootstrap.SetEnvConfig("test"); err != nil {
		panic(err)
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}

func setupTestDb(t *testing.T) (IConnection, string) {
	cfg := test.SetupTestCfg(t, test.REPO_REDIS)
	baseKey := cfg.Var.RedisKeyRoot

	conn := NewBaseDatabase(
		func() (master *redis.Client, slave *redis.Client, resultErr error) {
			connection, err := conn.Connect(cfg.ClubRedis)
			if err != nil {
				resultErr = err
				return
			}
			master = connection
			slave = connection
			return
		},
		func(connectionCreator IConnection) interface{} { return nil },
		baseKey,
	)
	t.Cleanup(func() {
		conn.Dispose()
	})
	return conn, baseKey
}

type testParser[Key any, Field any, Value any] struct {
	stringifyKey   func(key Key) string
	parseKey       func(keyStr string) (Key, error)
	stringifyField func(field Field) string
	parseField     func(fieldStr string) (Field, error)
	stringifyValue func(value Value) (string, error)
	parseValue     func(valueStr string) (Value, error)
}

func newTestParser[Key any, Field any, Value any](
	stringifyKey func(key Key) string,
	parseKey func(keyStr string) (Key, error),
	stringifyField func(field Field) string,
	parseField func(fieldStr string) (Field, error),
	stringifyValue func(value Value) (string, error),
	parseValue func(valueStr string) (Value, error),
) *testParser[Key, Field, Value] {
	return &testParser[Key, Field, Value]{
		stringifyKey:   stringifyKey,
		parseKey:       parseKey,
		stringifyField: stringifyField,
		parseField:     parseField,
		stringifyValue: stringifyValue,
		parseValue:     parseValue,
	}
}

func (t *testParser[Key, Field, Value]) set(
	p *testParser[Key, Field, Value],
) {
	if fn := p.stringifyKey; fn != nil {
		t.stringifyKey = fn
	}
	if fn := p.parseKey; fn != nil {
		t.parseKey = fn
	}
	if fn := p.stringifyField; fn != nil {
		t.stringifyField = fn
	}
	if fn := p.parseField; fn != nil {
		t.parseField = fn
	}
	if fn := p.stringifyValue; fn != nil {
		t.stringifyValue = fn
	}
	if fn := p.parseValue; fn != nil {
		t.parseValue = fn
	}
}

func (t testParser[Key, Field, Value]) StringifyKey(key Key) string {
	if fn := t.stringifyKey; fn != nil {
		return fn(key)
	}
	return ""
}

func (t testParser[Key, Field, Value]) ParseKey(keyStr string) (r Key, err error) {
	if fn := t.parseKey; fn != nil {
		return fn(keyStr)
	}
	err = fmt.Errorf("no implement")
	return
}

func (t testParser[Key, Field, Value]) StringifyField(field Field) string {
	if fn := t.stringifyField; fn != nil {
		return fn(field)
	}
	return ""
}

func (t testParser[Key, Field, Value]) ParseField(fieldStr string) (r Field, err error) {
	if fn := t.parseField; fn != nil {
		return fn(fieldStr)
	}
	err = fmt.Errorf("no implement")
	return
}

func (t testParser[Key, Field, Value]) StringifyValue(value Value) (string, error) {
	if fn := t.stringifyValue; fn != nil {
		return fn(value)
	}
	return "", fmt.Errorf("no implement")
}

func (t testParser[Key, Field, Value]) ParseValue(valueStr string) (r Value, err error) {
	if fn := t.parseValue; fn != nil {
		return fn(valueStr)
	}
	err = fmt.Errorf("no implement")
	return
}
