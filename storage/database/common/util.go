package common

import (
	"heroku-line-bot/global"
	"heroku-line-bot/storage/database/domain"
	"heroku-line-bot/util"
	"reflect"
	"time"
)

func NewLocalTime(t time.Time) domain.LocationTime {
	result := domain.LocationTime{}
	result.Scan(t)
	return result
}

func ConverTimeZone(dest interface{}) {
	destValue := reflect.ValueOf(dest)
	ConverTimeZoneValue(destValue)
}

func ConverTimeZoneValue(destValue reflect.Value) {
	k := destValue.Kind()
	switch k {
	case reflect.Ptr:
		ConverTimeZoneValue(destValue.Elem())
	case reflect.Array, reflect.Slice:
		len := destValue.Len()
		for i := 0; i < len; i++ {
			v := destValue.Index(i)
			ConverTimeZoneValue(v)
		}
	case reflect.Struct:
		if destValue.CanInterface() {
			destI := destValue.Interface()
			if t, ok := destI.(time.Time); ok {
				t = util.GetTimeIn(t, global.Location)
				newValue := reflect.ValueOf(t)
				destValue.Set(newValue)
				return
			}
		}

		destType := destValue.Type()
		for i := 0; i < destType.NumField(); i++ {
			v := destValue.Field(i)
			ConverTimeZoneValue(v)
		}
	}
}
