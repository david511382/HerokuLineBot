package common

import (
	"fmt"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/util"
	"strconv"
	"time"
)

type HourMinTime string

func NewHourMinTime(hour, min uint) HourMinTime {
	if hour >= 24 {
		hour = 0
	}
	if min >= 60 {
		min = 0
	}
	hourStr := strconv.Itoa(int(hour))
	minStr := strconv.Itoa(int(min))
	if l := len(hourStr); l < 2 {
		hourStr = "0" + hourStr
	} else if l > 2 {
		hourStr = hourStr[l-2 : l]
	}
	if l := len(minStr); l < 2 {
		minStr = "0" + minStr
	} else if l > 2 {
		minStr = minStr[l-2 : l]
	}
	return HourMinTime(fmt.Sprintf("%s:%s", hourStr, minStr))
}

func NewHourMinTimeOf(t time.Time) HourMinTime {
	return NewHourMinTime(uint(t.Hour()), uint(t.Minute()))
}

func (t HourMinTime) ToString() string {
	return string(t)
}

func (t HourMinTime) Time() (resultTime time.Time, resultErr error) {
	rawT, err := time.Parse("15:04", t.ToString())
	if err != nil {
		resultErr = err
		return
	}
	tp := util.GetTimePLoc(global.Location, 0, 1, 1, rawT.Hour(), rawT.Minute())
	resultTime = *tp
	return
}

func (t HourMinTime) ForceTime() time.Time {
	resultTime, _ := t.Time()
	return resultTime
}

type MinSecTime string

func NewMinSecTime(min, sec uint) MinSecTime {
	if min >= 60 {
		min = 0
	}
	if sec >= 60 {
		sec = 0
	}
	minStr := strconv.Itoa(int(min))
	secStr := strconv.Itoa(int(sec))
	if l := len(minStr); l < 2 {
		minStr = "0" + minStr
	} else if l > 2 {
		minStr = minStr[l-2 : l]
	}
	if l := len(secStr); l < 2 {
		secStr = "0" + secStr
	} else if l > 2 {
		secStr = secStr[l-2 : l]
	}
	return MinSecTime(fmt.Sprintf("%s:%s", minStr, secStr))
}

func (t MinSecTime) ToString() string {
	return string(t)
}

func (t MinSecTime) Time() (resultTime time.Time, resultErr error) {
	rawT, err := time.Parse("04:05", t.ToString())
	if err != nil {
		resultErr = err
		return
	}
	tp := util.GetTimePLoc(global.Location, 0, 1, 1, 0, rawT.Minute(), rawT.Second())
	resultTime = *tp
	return
}

func (t MinSecTime) ForceTime() time.Time {
	resultTime, _ := t.Time()
	return resultTime
}
