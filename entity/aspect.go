package entity

import (
	"fmt"
	"github.com/jucardi/go-db/logger"
	"reflect"
)

const (
	MethodBeforeCreate = "BeforeCreate"
	MethodBeforeUpdate = "BeforeUpdate"
	MethodBeforeDelete = "BeforeDelete"

	MethodAfterCreate = "AfterCreate"
	MethodAfterUpdate = "AfterUpdate"
	MethodAfterDelete = "AfterDelete"
	MethodAfterFound  = "AfterFound"
)

var errType = reflect.TypeOf((*error)(nil))

func Invoke(method string, entity ...interface{}) error {
	for _, e := range entity {
		val := reflect.ValueOf(e)
		if IsNil(e) {
			logger.Get().Warn(fmt.Sprintf("failed to invoke %s, entity appears to be nil", method))
			continue
		}
		m := val.MethodByName(method)
		if !m.IsValid() && val.Kind() == reflect.Ptr {
			m = val.Elem().MethodByName(method)
		}
		if !m.IsValid() {
			logger.Get().Debug(method, " not found")
			continue
		}
		if m.Kind() != reflect.Func {
			return fmt.Errorf("failed to invoke %s, is not a function", method)
		}
		rets := m.Call(nil)
		if len(rets) != 1 {
			return fmt.Errorf("failed to invoke %s, expected 1 return value, found %d", method, len(rets))
		}
		if ret, ok := rets[0].Interface().(error); ok && ret != nil {
			return fmt.Errorf("failed to invoke %s, %s", method, ret)
		} else if !rets[0].Type().Implements(errType) {
			return fmt.Errorf("failed to invoke %s, incorrect return type, expected (error)", method)
		}
	}
	return nil
}
