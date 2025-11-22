package print

import "fmt"

func Entry(label, value string) {
	if value == "" {
		fmt.Printf("  %s: \"\"\n", label)
	} else {
		fmt.Printf("  %s: %s\n", label, value)
	}
}
