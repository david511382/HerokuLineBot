package common

import (
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/util"
	"time"
)

func TimeInt(t time.Time, tt util.TimeType) int {
	t = t.In(global.Location)
	return util.TimeInt(t, tt)
}

func ClockInt(t time.Time, tt util.TimeType) int {
	t = t.In(global.Location)
	return util.ClockInt(t, tt)
}

func GetTime(ts ...int) time.Time {
	return *GetTimeP(ts...)
}

func GetTimeP(ts ...int) *time.Time {
	return util.GetTimePLoc(global.Location, ts...)
}
