package util

import (
	"strconv"
)

// StringToInt converts a string to an integer
// Returns the integer value and an error if conversion fails
func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}
