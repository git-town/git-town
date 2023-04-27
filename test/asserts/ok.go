package asserts

import "fmt"

func Ok(err error, args ...any) {
	if err == nil {
		return
	}
	if len(args) == 0 {
		panic(err)
	}
	first, ok := args[0].(string)
	if !ok {
		panic("please provide a string as the formatting argument")
	}
	panic(fmt.Sprintf(first, args[1:]))
}
