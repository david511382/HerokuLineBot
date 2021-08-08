package domain

import (
	"database/sql/driver"
	"heroku-line-bot/util"
	"time"

	commonLogic "heroku-line-bot/logic/common"
	errLogic "heroku-line-bot/logic/error"
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
		t.time = util.GetTimeIn(tt, commonLogic.Location)
	} else {
		return errLogic.Newf("can not convert %v to time", v)
	}

	return nil
}
