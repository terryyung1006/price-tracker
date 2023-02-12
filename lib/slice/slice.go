package slice

import "fmt"

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func GetSliceOfKeys[T comparable, V interface{}](targetMap map[T]V) []T {

	keys := make([]T, 0, len(targetMap))
	for k := range targetMap {
		keys = append(keys, k)
	}
	return keys
}

func MinInt(array []int) (int, error) {
	if len(array) == 0 {
		return -1, fmt.Errorf("[MinInt] len of array cannot be 0")
	}
	var min int = array[0]
	for _, value := range array {
		if min > value {
			min = value
		}
	}
	return min, nil
}

func IndexOfMinInt(array []int) (int, error) {
	if len(array) == 0 {
		return -1, fmt.Errorf("[IndexOfMinInt] len of array cannot be 0")
	}
	var min int = array[0]
	var index int = 0
	for i, value := range array {
		if min > value {
			min = value
			index = i
		}
	}
	return index, nil
}
