package util

import (
	"log"
	"os"
)

// DoesFileExist returns whether or not a file exists at the given path
func DoesFileExist(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatal("Error getting stat for run result file:", err)
	}
	return true
}
