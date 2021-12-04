package util

import "reflect"

type IFindCondition interface {
	HandleFoundValue(foundValue *reflect.Value) (isFound bool)
}

func ReflectFindProperty(dest interface{}, condition IFindCondition) {
	if condition == nil {
		return
	}

	destValue := reflect.ValueOf(dest)
	reflectFindProperty(destValue, condition)
}

func reflectFindProperty(destValue reflect.Value, condition IFindCondition) {
	if condition == nil {
		return
	}

	k := destValue.Kind()
	switch k {
	case reflect.Ptr:
		reflectFindProperty(destValue.Elem(), condition)
	case reflect.Array, reflect.Slice:
		len := destValue.Len()
		for i := 0; i < len; i++ {
			v := destValue.Index(i)
			reflectFindProperty(v, condition)
		}
	case reflect.Struct:
		if destValue.CanSet() && destValue.CanInterface() {
			if ok := condition.HandleFoundValue(&destValue); ok {
				return
			}
		}

		destType := destValue.Type()
		for i := 0; i < destType.NumField(); i++ {
			v := destValue.Field(i)
			reflectFindProperty(v, condition)
		}
	}
}
