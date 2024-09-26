package utils

import (
	"strconv"
)

func StringToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// handle error appropriately, e.g., log it or set a default value
		return 0.0
	}
	return f
}
