package numbers

import "strconv"

// 将字符串转换成int
// @param str 字符串
func StringToIntMust(str string) int {
	atoi, _ := strconv.Atoi(str)
	return atoi
}

// 见字符串转换成int64
// @param str 字符串
func StringToInt64Must(str string) int64 {
	atoi, _ := strconv.ParseInt(str, 10, 64)
	return atoi
}
