package entity

import "reflect"

func IsNil(obj interface{}) bool {
	val, ok := obj.(reflect.Value)
	if !ok {
		val = reflect.ValueOf(obj)
	}
	return obj == nil || !val.IsValid() || (val.Kind() == reflect.Ptr && val.IsNil())
}
