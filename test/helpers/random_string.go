package helpers

import (
	"math/rand"
	"strconv"
)

// RandomString provides a string containing the given amount of random numbers.
func RandomString(length int) (result string) {
	for i := 0; i < length; i++ {
		result += strconv.Itoa(rand.Intn(10))
	}
	return result
}
