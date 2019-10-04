package helpers

import (
	"math/rand"
	"strconv"
)

// RandomNumberString provides a string containing the given amount of random numbers.
func RandomNumberString(length int) (result string) {
	for i := 0; i < length; i++ {
		result += strconv.Itoa(rand.Intn(10))
	}
	return result
}
