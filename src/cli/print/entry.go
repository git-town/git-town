package print

import "fmt"

func Entry(label, value string) {
	fmt.Printf("  %s: %s\n", label, value)
}
