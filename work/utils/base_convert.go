package utils

import (
	"math"
	"strconv"
	"strings"
)

//小写字符去除0 l o
var tenTo33 map[int]string = map[int]string{0: "z", 1: "1", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6", 7: "7", 8: "8", 9: "9", 10: "a", 11: "b", 12: "c", 13: "d", 14: "e", 15: "f", 16: "g", 17: "h", 18: "i", 19: "j", 20: "k", 21: "m", 22: "n", 23: "p", 24: "q", 25: "r", 26: "s", 27: "t", 28: "u", 29: "v", 30: "w", 31: "x", 32: "y"}

//大小写字符去除0 O o l I
var tenTo57 map[int]string = map[int]string{0: "z", 1: "1", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6", 7: "7", 8: "8", 9: "9", 10: "a", 11: "b", 12: "c", 13: "d", 14: "e", 15: "f", 16: "g", 17: "h", 18: "i", 19: "j", 20: "k", 21: "m", 22: "n", 23: "p", 24: "q", 25: "r", 26: "s", 27: "t", 28: "u", 29: "v", 30: "w", 31: "x", 32: "y", 33: "A", 34: "B", 35: "C", 36: "D", 37: "E", 38: "F", 39: "G", 40: "H", 41: "J", 42: "K", 43: "L", 44: "M", 45: "N", 46: "P", 47: "Q", 48: "R", 49: "S", 50: "T", 51: "U", 52: "V", 53: "W", 54: "X", 55: "Y", 56: "Z"}

var tenTo76 map[int]string = map[int]string{0: "0", 1: "1", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6", 7: "7", 8: "8", 9: "9", 10: "a", 11: "b", 12: "c", 13: "d", 14: "e", 15: "f", 16: "g", 17: "h", 18: "i", 19: "j", 20: "k", 21: "l", 22: "m", 23: "n", 24: "o", 25: "p", 26: "q", 27: "r", 28: "s", 29: "t", 30: "u", 31: "v", 32: "w", 33: "x", 34: "y", 35: "z", 36: ":", 37: ";", 38: "<", 39: "=", 40: ">", 41: "?", 42: "@", 43: "[", 44: "]", 45: "^", 46: "_", 47: "{", 48: "|", 49: "}", 50: "A", 51: "B", 52: "C", 53: "D", 54: "E", 55: "F", 56: "G", 57: "H", 58: "I", 59: "J", 60: "K", 61: "L", 62: "M", 63: "N", 64: "O", 65: "P", 66: "Q", 67: "R", 68: "S", 69: "T", 70: "U", 71: "V", 72: "W", 73: "X", 74: "Y", 75: "Z"}

//==============================================最多76进制=====================================================

// 10进制转最多76进制
func decimalToAny(num, n int, base map[int]string) string {
	newNumStr := ""
	var remainder int
	var remainderString string
	baseLen := len(base)
	for num != 0 {
		remainder = num % n
		isSearch := false
		switch baseLen {
		case 33:
			isSearch = remainder == 0 || (33 > remainder && remainder > 9)
		case 57:
			isSearch = remainder == 0 || (57 > remainder && remainder > 9)
		case 76:
			isSearch = 76 > remainder && remainder > 9
		}
		if isSearch {
			remainderString = base[remainder]
		} else {
			remainderString = strconv.Itoa(remainder)
		}
		newNumStr = remainderString + newNumStr
		num = num / n
	}
	return newNumStr
}

// 最多76进制转10进制
func anyToDecimal(num string, n int, base map[int]string) int {
	var newNum float64
	newNum = 0.0
	nNum := len(strings.Split(num, "")) - 1
	for _, value := range strings.Split(num, "") {
		tmp := float64(findKey(value, base))
		if tmp != -1 {
			newNum = newNum + tmp*math.Pow(float64(n), float64(nNum))
			nNum = nNum - 1
		} else {
			break
		}
	}
	return int(newNum)
}

// map根据value找key
func findKey(in string, base map[int]string) int {
	result := -1
	for k, v := range base {
		if in == v {
			result = k
		}
	}
	return result
}
