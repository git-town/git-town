package gohacks

import "reflect"

// TypeName provides the name of the type of the given variable.
func TypeName(myvar interface{}) string {
	t := reflect.TypeOf(myvar)
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}
	return t.Name()
}
