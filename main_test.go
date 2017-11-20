// See ./test-coverage.md for why this file exists

package main

import (
	"os"
	"testing"
)

func TestRunMain(t *testing.T) {

	// delete the coverage measure parameter
	os.Args = append(os.Args[:1], os.Args[2:]...)

	main()
}
