package common

import (
	"fmt"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"testing"
)

func TestBaseHashKey_HSet(t *testing.T) {
	t.Parallel()

	type args struct {
		field string
		value string
	}
	type migrations struct {
		redisData map[string]string
		parser    *testParser[string, string, string]
	}
	type wants struct {
		redisData map[string]string
		errInfo   errUtil.IError
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"NOT_CHANGE",
			args{
				"k",
				"v",
			},
			migrations{
				redisData: map[string]string{
					"k": "v",
				},
			},
			wants{
				redisData: map[string]string{
					"k": "v",
				},
				errInfo: NotChangeErrInfo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connection, baseKey := setupTestDb(t)

			parser := *defaultParser
			if p := tt.migrations.parser; p != nil {
				parser.set(p)
			}
			key := NewBaseHashKey[string, string](
				connection,
				baseKey,
				parser,
			)

			migrationKey := NewBaseHashKey[string, string](
				connection,
				baseKey,
				defaultParser,
			)
			if errInfo := migrationKey.Migration(tt.migrations.redisData); errInfo != nil {
				t.Fatal(errInfo.Error())
			}

			errInfo := key.HSet(tt.args.field, tt.args.value)
			if got, want := errInfo, tt.wants.errInfo; !errUtil.Equal(got, want) {
				if got != nil {
					t.Error(got.Error())
				}
				if want != nil {
					t.Error(want.Error())
				}
			}

			if got, err := key.Read(); err != nil {
				t.Fatal(err.Error())
			} else {
				if ok, msg := util.Comp(got, tt.wants.redisData); !ok {
					t.Errorf(msg)
					return
				}
			}
		})
	}
}

func TestBaseHashKey_HMSet(t *testing.T) {
	t.Parallel()

	type args struct {
		fieldValueMap map[string]string
	}
	type migrations struct {
		redisData map[string]string
		parser    *testParser[string, string, string]
	}
	type wants struct {
		redisData map[string]string
		errInfo   errUtil.IError
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"fail value parse",
			args{
				map[string]string{
					"f": "f",
					"k": "v",
				},
			},
			migrations{
				parser: newTestParser[string, string](
					nil, nil, nil, nil,
					func(value string) (string, error) {
						if value == "f" {
							return "", fmt.Errorf("fail")
						}
						return value, nil
					},
					nil,
				),
			},
			wants{
				redisData: map[string]string{},
				errInfo:   errUtil.New("fail"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connection, baseKey := setupTestDb(t)

			parser := *defaultParser
			if p := tt.migrations.parser; p != nil {
				parser.set(p)
			}
			key := NewBaseHashKey[string, string](
				connection,
				baseKey,
				parser,
			)

			migrationKey := NewBaseHashKey[string, string](
				connection,
				baseKey,
				defaultParser,
			)
			if errInfo := migrationKey.Migration(tt.migrations.redisData); errInfo != nil {
				t.Fatal(errInfo.Error())
			}

			errInfo := key.HMSet(tt.args.fieldValueMap)
			if got, want := errInfo, tt.wants.errInfo; !errUtil.Equal(got, want) {
				if got != nil {
					t.Error(got.Error())
				}
				if want != nil {
					t.Error(want.Error())
				}
			}

			if got, err := key.Read(); err != nil {
				t.Fatal(err.Error())
			} else {
				if ok, msg := util.Comp(got, tt.wants.redisData); !ok {
					t.Errorf(msg)
					return
				}
			}
		})
	}
}

func TestBaseHashKey_HKeys(t *testing.T) {
	t.Parallel()

	type args struct {
	}
	type migrations struct {
		redisData map[string]string
		parser    *testParser[string, string, string]
	}
	type wants struct {
		fields  []string
		isError bool
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"fail field parse",
			args{},
			migrations{
				redisData: map[string]string{
					"k": "v",
					"f": "f",
				},
				parser: newTestParser[string, string, string](
					nil, nil, nil,
					func(fieldStr string) (string, error) {
						if fieldStr == "f" {
							return "", fmt.Errorf("fail")
						}
						return fieldStr, nil
					},
					nil,
					nil,
				),
			},
			wants{
				isError: true,
			},
		},
		{
			"no key",
			args{},
			migrations{
				redisData: map[string]string{},
			},
			wants{
				fields: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connection, baseKey := setupTestDb(t)

			parser := *defaultParser
			if p := tt.migrations.parser; p != nil {
				parser.set(p)
			}
			key := NewBaseHashKey[string, string](
				connection,
				baseKey,
				parser,
			)

			migrationKey := NewBaseHashKey[string, string](
				connection,
				baseKey,
				defaultParser,
			)
			if errInfo := migrationKey.Migration(tt.migrations.redisData); errInfo != nil {
				t.Fatal(errInfo.Error())
			}

			fields, errInfo := key.HKeys()
			if errInfo != nil && !tt.wants.isError ||
				tt.wants.isError != (errInfo != nil && errInfo.IsError()) {
				if errInfo != nil {
					t.Error(errInfo.Error())
				} else {
					t.Error("want error")
				}
				return
			} else if tt.wants.isError {
				return
			}

			if ok, msg := util.Comp(fields, tt.wants.fields); !ok {
				t.Errorf(msg)
				return
			}
		})
	}
}

func TestBaseHashKey_HGetAll(t *testing.T) {
	t.Parallel()

	type args struct {
	}
	type migrations struct {
		redisData map[string]string
		parser    *testParser[string, string, string]
	}
	type wants struct {
		redisData map[string]string
		isError   bool
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"fail field parse",
			args{},
			migrations{
				redisData: map[string]string{
					"k": "v",
					"f": "v",
				},
				parser: newTestParser[string, string, string](
					nil, nil, nil,
					func(fieldStr string) (string, error) {
						if fieldStr == "f" {
							return "", fmt.Errorf("fail")
						}
						return fieldStr, nil
					},
					nil,
					nil,
				),
			},
			wants{
				isError: true,
			},
		},
		{
			"fail value parse",
			args{},
			migrations{
				redisData: map[string]string{
					"k": "v",
					"v": "f",
				},
				parser: newTestParser[string, string](
					nil, nil, nil, nil, nil,
					func(valueStr string) (string, error) {
						if valueStr == "f" {
							return "", fmt.Errorf("fail")
						}
						return valueStr, nil
					},
				),
			},
			wants{
				isError: true,
			},
		},
		{
			"no key",
			args{},
			migrations{
				redisData: map[string]string{},
			},
			wants{
				redisData: map[string]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connection, baseKey := setupTestDb(t)

			parser := *defaultParser
			if p := tt.migrations.parser; p != nil {
				parser.set(p)
			}
			key := NewBaseHashKey[string, string](
				connection,
				baseKey,
				parser,
			)

			migrationKey := NewBaseHashKey[string, string](
				connection,
				baseKey,
				defaultParser,
			)
			if errInfo := migrationKey.Migration(tt.migrations.redisData); errInfo != nil {
				t.Fatal(errInfo.Error())
			}

			redisData, errInfo := key.HGetAll()
			if errInfo != nil && !tt.wants.isError ||
				tt.wants.isError != (errInfo != nil && errInfo.IsError()) {
				if errInfo != nil {
					t.Error(errInfo.Error())
				} else {
					t.Error("want error")
				}
				return
			} else if tt.wants.isError {
				return
			}

			if ok, msg := util.Comp(redisData, tt.wants.redisData); !ok {
				t.Errorf(msg)
				return
			}
		})
	}
}

func TestBaseHashKey_HGet(t *testing.T) {
	t.Parallel()

	type args struct {
		field string
	}
	type migrations struct {
		redisData map[string]string
		parser    *testParser[string, string, string]
	}
	type wants struct {
		value   *string
		errInfo errUtil.IError
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"fail value parse",
			args{
				field: "k",
			},
			migrations{
				redisData: map[string]string{
					"k": "f",
				},
				parser: newTestParser[string, string](
					nil, nil, nil, nil, nil,
					func(valueStr string) (string, error) {
						if valueStr == "f" {
							return "", fmt.Errorf("fail")
						}
						return valueStr, nil
					},
				),
			},
			wants{
				errInfo: errUtil.New("fail"),
			},
		},
		{
			"NOT_EXIST",
			args{
				field: "k",
			},
			migrations{
				redisData: map[string]string{},
			},
			wants{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connection, baseKey := setupTestDb(t)

			parser := *defaultParser
			if p := tt.migrations.parser; p != nil {
				parser.set(p)
			}
			key := NewBaseHashKey[string, string](
				connection,
				baseKey,
				parser,
			)

			migrationKey := NewBaseHashKey[string, string](
				connection,
				baseKey,
				defaultParser,
			)
			if errInfo := migrationKey.Migration(tt.migrations.redisData); errInfo != nil {
				t.Fatal(errInfo.Error())
			}

			redisData, errInfo := key.HGet(tt.args.field)
			if got, want := errInfo, tt.wants.errInfo; !errUtil.Equal(got, want) {
				if got != nil {
					t.Error(got.Error())
				}
				if want != nil {
					t.Error(want.Error())
				}
			}

			if ok, msg := util.Comp(redisData, tt.wants.value); !ok {
				t.Errorf(msg)
				return
			}
		})
	}
}

func TestBaseHashKey_HMGet(t *testing.T) {
	t.Parallel()

	type args struct {
		fields []string
	}
	type migrations struct {
		redisData map[string]string
		parser    *testParser[string, string, string]
	}
	type wants struct {
		redisData map[string]string
		errInfo   errUtil.IError
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"fail value parse",
			args{
				fields: []string{"k", "k1"},
			},
			migrations{
				redisData: map[string]string{
					"k":  "v",
					"k1": "f",
				},
				parser: newTestParser[string, string](
					nil, nil, nil, nil, nil,
					func(valueStr string) (string, error) {
						if valueStr == "f" {
							return "", fmt.Errorf("fail")
						}
						return valueStr, nil
					},
				),
			},
			wants{
				errInfo: func() errUtil.IError {
					errInfo := errUtil.New("fail")
					errInfo.Attr("value", "f")
					return errInfo
				}(),
			},
		},
		{
			"no key",
			args{
				fields: []string{"k"},
			},
			migrations{
				redisData: map[string]string{
					"k1": "f",
				},
			},
			wants{
				redisData: map[string]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connection, baseKey := setupTestDb(t)

			parser := *defaultParser
			if p := tt.migrations.parser; p != nil {
				parser.set(p)
			}
			key := NewBaseHashKey[string, string](
				connection,
				baseKey,
				parser,
			)

			migrationKey := NewBaseHashKey[string, string](
				connection,
				baseKey,
				defaultParser,
			)
			if errInfo := migrationKey.Migration(tt.migrations.redisData); errInfo != nil {
				t.Fatal(errInfo.Error())
			}

			redisData, errInfo := key.HMGet(tt.args.fields...)
			if got, want := errInfo, tt.wants.errInfo; !errUtil.Equal(got, want) {
				if got != nil {
					t.Error(got.Error())
				}
				if want != nil {
					t.Error(want.Error())
				}
			}
			if e := tt.wants.errInfo; e != nil && e.IsError() {
				return
			}

			if ok, msg := util.Comp(redisData, tt.wants.redisData); !ok {
				t.Errorf(msg)
				return
			}
		})
	}
}

func TestBaseHashKey_HDel(t *testing.T) {
	t.Parallel()

	type args struct {
		fields []string
	}
	type migrations struct {
		redisData map[string]string
		parser    *testParser[string, string, string]
	}
	type wants struct {
		redisData map[string]string
		change    int64
		errInfo   errUtil.IError
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"no key",
			args{
				fields: []string{"k"},
			},
			migrations{
				redisData: map[string]string{
					"k1": "f",
				},
			},
			wants{
				redisData: map[string]string{
					"k1": "f",
				},
				change: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connection, baseKey := setupTestDb(t)

			parser := *defaultParser
			if p := tt.migrations.parser; p != nil {
				parser.set(p)
			}
			key := NewBaseHashKey[string, string](
				connection,
				baseKey,
				parser,
			)

			migrationKey := NewBaseHashKey[string, string](
				connection,
				baseKey,
				defaultParser,
			)
			if errInfo := migrationKey.Migration(tt.migrations.redisData); errInfo != nil {
				t.Fatal(errInfo.Error())
			}

			change, errInfo := key.HDel(tt.args.fields...)
			if got, want := errInfo, tt.wants.errInfo; !errUtil.Equal(got, want) {
				if got != nil {
					t.Error(got.Error())
				}
				if want != nil {
					t.Error(want.Error())
				}
			}
			if e := tt.wants.errInfo; e != nil && e.IsError() {
				return
			}
			if ok, msg := util.Comp(change, tt.wants.change); !ok {
				t.Errorf(msg)
				return
			}

			if got, err := key.Read(); err != nil {
				t.Fatal(err.Error())
			} else {
				if ok, msg := util.Comp(got, tt.wants.redisData); !ok {
					t.Errorf(msg)
					return
				}
			}
		})
	}
}
