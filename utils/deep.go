package utils

import "reflect"

func DeepClone[V1 any](v V1) V1 {
	val := reflect.ValueOf(v)
	if val.IsNil() {
		return v
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	result := reflect.New(val.Type()).Elem()
	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldType := field.Type()
			if field.CanAddr() {
				addr := field.Addr()
				if addr.Kind() == reflect.Ptr && !addr.IsNil() {
					cloned := DeepClone(addr.Interface())
					result.Field(i).Set(reflect.ValueOf(cloned))
				} else {
					result.Field(i).Set(field)
				}
			} else if field.Kind() == reflect.Slice {
				s := reflect.MakeSlice(fieldType, field.Len(), field.Cap())
				reflect.Copy(s, field)
				result.Field(i).Set(s)
			} else if field.Kind() == reflect.Map {
				m := reflect.MakeMap(fieldType)
				keys := field.MapKeys()
				for _, key := range keys {
					m.SetMapIndex(key, field.MapIndex(key))
				}
				result.Field(i).Set(m)
			} else {
				result.Field(i).Set(field)
			}
		}
	case reflect.Slice:
		s := reflect.MakeSlice(val.Type(), val.Len(), val.Cap())
		reflect.Copy(s, val)
		return s.Interface().(V1)
	case reflect.Map:
		m := reflect.MakeMap(val.Type())
		keys := val.MapKeys()
		for _, key := range keys {
			m.SetMapIndex(key, val.MapIndex(key))
		}
		return m.Interface().(V1)
	}
	return result.Interface().(V1)
}
