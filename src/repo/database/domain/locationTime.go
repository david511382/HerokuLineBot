package domain

import (
	"database/sql/driver"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"time"
)

type LocationTime struct {
	time time.Time
}

func (t *LocationTime) Time() time.Time {
	return t.time
}

func (t *LocationTime) TimeP() *time.Time {
	if t == nil {
		return nil
	}
	return &t.time
}

func (t LocationTime) Value() (driver.Value, error) {
	if t.time.IsZero() {
		return nil, nil
	}
	return t.time, nil
}

func (t *LocationTime) Scan(v interface{}) error {
	if tt, ok := v.(time.Time); ok {
		t.time = util.GetTimeIn(tt, global.TimeUtilObj.GetLocation())
	} else {
		return errUtil.Newf("can not convert %v to time", v)
	}

	return nil
}
