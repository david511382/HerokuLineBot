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
	destValue := reflect.ValueOf(dest)
	l.ConvertReflect(destValue)
}

func (l LocationConverter) ConvertReflect(destValue reflect.Value) {
	k := destValue.Kind()
	switch k {
	case reflect.Ptr:
		l.ConvertReflect(destValue.Elem())
	case reflect.Array, reflect.Slice:
		len := destValue.Len()
		for i := 0; i < len; i++ {
			v := destValue.Index(i)
			l.ConvertReflect(v)
		}
	case reflect.Struct:
		if destValue.CanSet() && destValue.CanInterface() {
			destI := destValue.Interface()
			t, ok := destI.(time.Time)
			if ok {
				var newValue reflect.Value
				if l.isChangeLocationDirectly {
					t = GetTimeIn(t, l.location)
					newValue = reflect.ValueOf(t)
				} else {
					t = t.In(l.location)
					newValue = reflect.ValueOf(t)
				}

				destValue.Set(newValue)
				return
			}
		}

		destType := destValue.Type()
		for i := 0; i < destType.NumField(); i++ {
			v := destValue.Field(i)
			l.ConvertReflect(v)
		}
	}
}
