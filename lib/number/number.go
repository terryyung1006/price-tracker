package number

import (
	"fmt"
	"strconv"
)

func MiliSecondToSecondTimestamp(timestamp int) (int, error) {
	numOfDigit := len(strconv.Itoa(timestamp))
	if numOfDigit == 10 {
		return timestamp, nil
	}
	if numOfDigit == 13 {
		return timestamp / 1000, nil
	}
	return -1, fmt.Errorf("[MiliSecondToSecondTimestamp] input timestamp num of digits invalid, num of digits: %v", numOfDigit)
}

func Abs(val int) int {
	if val < 0 {
		val *= -1
	}
	return val
}
