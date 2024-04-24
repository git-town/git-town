package gohacks

import "reflect"

// IsNil detects nil even when evaluating pointer interface variables.
func IsNil(val any) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	k := v.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer,
		reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	}
	return false
}
