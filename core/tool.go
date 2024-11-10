package core

import "strconv"

func ToFloat64(str string) float64 {
	floatVal, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(err)
	}
	return floatVal
}
