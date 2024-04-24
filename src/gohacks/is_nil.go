package gohacks

import "reflect"

// IsNil detects nil even when evaluating pointer interface variables.
func IsNil(val any) bool {
	if val == nil {
		return true
	}
	value := reflect.ValueOf(val)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer,
		reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return value.IsNil()
	}
	return false
}
