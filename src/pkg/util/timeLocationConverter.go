package util

import (
	"reflect"
	"time"
)

type LocationConverter struct {
	location                 *time.Location
	isChangeLocationDirectly bool
}

// isChangeLocationDirectly:true change location directly
func NewLocationConverter(location *time.Location, isChangeLocationDirectly bool) LocationConverter {
	return LocationConverter{
		location:                 location,
		isChangeLocationDirectly: isChangeLocationDirectly,
	}
}

func (l LocationConverter) GetTime(ts ...int) time.Time {
	return *l.GetTimeP(ts...)
}

func (l LocationConverter) GetTimeP(ts ...int) *time.Time {
	return GetTimePLoc(l.location, ts...)
}

func (l LocationConverter) ConvertTime(t time.Time) time.Time {
	if l.isChangeLocationDirectly {
		return GetTimeIn(t, l.location)
	} else {
		return t.In(l.location)
	}

}

func (l LocationConverter) Convert(dest interface{}) {
	ReflectFindProperty(dest, l)
}

func (l LocationConverter) HandleFoundValue(foundValue *reflect.Value) (isFound bool) {
	destI := foundValue.Interface()
	t, ok := destI.(time.Time)
	isFound = ok
	if ok {
		var newValue reflect.Value
		if l.isChangeLocationDirectly {
			t = GetTimeIn(t, l.location)
			newValue = reflect.ValueOf(t)
		} else {
			t = t.In(l.location)
			newValue = reflect.ValueOf(t)
		}

		foundValue.Set(newValue)
		return
	}
	return
}
