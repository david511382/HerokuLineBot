package common

import (
	"fmt"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"testing"
	"time"
)

func TestBaseKey_SetNX(t *testing.T) {
	t.Parallel()

	type args struct {
		value string
		et    time.Duration
	}
	type migrations struct {
		redisData *string
		parser    *testParser[string, string, string]
	}
	type wants struct {
		redisData        *string
		timeoutRedisData *string
		errInfo          errUtil.IError
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
				value: "v",
				et:    0,
			},
			migrations{
				parser: newTestParser[string, string](
					nil, nil, nil, nil,
					func(value string) (string, error) {
						if value == "v" {
							return "", fmt.Errorf("fail")
						}
						return value, nil
					},
					nil,
				),
			},
			wants{
				errInfo: errUtil.New("fail"),
			},
		},
		{
			"timeout",
			args{
				value: "v",
				et:    time.Millisecond * 100,
			},
			migrations{
				redisData: nil,
			},
			wants{
				redisData:        util.PointerOf("v"),
				timeoutRedisData: nil,
			},
		},
		{
			"NOT_CHANGE",
			args{
				value: "v",
				et:    0,
			},
			migrations{
				redisData: util.PointerOf("v"),
			},
			wants{
				redisData:        util.PointerOf("v"),
				timeoutRedisData: util.PointerOf("v"),
				errInfo:          NotChangeErrInfo,
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
			key := NewBaseKey[string](
				connection,
				baseKey,
				parser,
			)

			migrationKey := NewBaseKey[string](
				connection,
				baseKey,
				defaultParser,
			)
			if p := tt.migrations.redisData; p != nil {
				if errInfo := migrationKey.Migration(*p); errInfo != nil {
					t.Fatal(errInfo.Error())
				}
			}

			errInfo := key.SetNX(tt.args.value, tt.args.et)
			if got, want := errInfo, tt.wants.errInfo; !errUtil.Equal(got, want) {
				if got != nil {
					t.Error(got.Error())
				}
				if want != nil {
					t.Error(want.Error())
				}
			}
			if errInfo != nil && errInfo.IsError() {
				return
			}

			if got, errInfo := key.Get(); errInfo != nil {
				t.Fatal(errInfo.Error())
			} else {
				if ok, msg := util.Comp(got, tt.wants.redisData); !ok {
					t.Errorf(msg)
					return
				}
			}

			time.Sleep(tt.args.et)

			{
				got, errInfo := key.Get()
				if errInfo != nil && errInfo.IsError() {
					t.Fatal(errInfo.Error())
				} else {
					if ok, msg := util.Comp(got, tt.wants.timeoutRedisData); !ok {
						t.Errorf(msg)
						return
					}
				}
			}
		})
	}
}

func TestBaseKey_Set(t *testing.T) {
	t.Parallel()

	type args struct {
		value string
		et    time.Duration
	}
	type migrations struct {
		redisData *string
		parser    *testParser[string, string, string]
	}
	type wants struct {
		redisData *string
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
				value: "v",
				et:    0,
			},
			migrations{
				parser: newTestParser[string, string](
					nil, nil, nil, nil,
					func(value string) (string, error) {
						if value == "v" {
							return "", fmt.Errorf("fail")
						}
						return value, nil
					},
					nil,
				),
			},
			wants{
				errInfo: errUtil.New("fail"),
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
			key := NewBaseKey[string](
				connection,
				baseKey,
				parser,
			)

			migrationKey := NewBaseKey[string](
				connection,
				baseKey,
				defaultParser,
			)
			if p := tt.migrations.redisData; p != nil {
				if errInfo := migrationKey.Migration(*p); errInfo != nil {
					t.Fatal(errInfo.Error())
				}
			}

			errInfo := key.Set(tt.args.value, tt.args.et)
			if got, want := errInfo, tt.wants.errInfo; !errUtil.Equal(got, want) {
				if got != nil {
					t.Error(got.Error())
				}
				if want != nil {
					t.Error(want.Error())
				}
			}
			if errInfo != nil && errInfo.IsError() {
				return
			}

			if got, errInfo := key.Get(); errInfo != nil {
				t.Fatal(errInfo.Error())
			} else {
				if ok, msg := util.Comp(got, tt.wants.redisData); !ok {
					t.Errorf(msg)
					return
				}
			}
		})
	}
}

func TestBaseKey_Get(t *testing.T) {
	t.Parallel()

	type args struct{}
	type migrations struct {
		redisData *string
		parser    *testParser[string, string, string]
	}
	type wants struct {
		redisData *string
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
			args{},
			migrations{
				redisData: util.PointerOf("v"),
				parser: newTestParser[string, string](
					nil, nil, nil, nil, nil,
					func(value string) (string, error) {
						if value == "v" {
							return "", fmt.Errorf("fail")
						}
						return value, nil
					},
				),
			},
			wants{
				errInfo: errUtil.New("fail"),
			},
		},
		{
			"NOT_EXIST",
			args{},
			migrations{},
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
			key := NewBaseKey[string](
				connection,
				baseKey,
				parser,
			)

			migrationKey := NewBaseKey[string](
				connection,
				baseKey,
				defaultParser,
			)
			if p := tt.migrations.redisData; p != nil {
				if errInfo := migrationKey.Migration(*p); errInfo != nil {
					t.Fatal(errInfo.Error())
				}
			}

			got, errInfo := key.Get()
			if got, want := errInfo, tt.wants.errInfo; !errUtil.Equal(got, want) {
				if got != nil {
					t.Error(got.Error())
				}
				if want != nil {
					t.Error(want.Error())
				}
			}
			if errInfo != nil && errInfo.IsError() {
				return
			}

			if ok, msg := util.Comp(got, tt.wants.redisData); !ok {
				t.Errorf(msg)
				return
			}
		})
	}
}
