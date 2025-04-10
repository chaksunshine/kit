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

// 将字符串转换成Int32
func StringToInt32Must(str string) int32 {
	atoi, _ := strconv.ParseInt(str, 10, 32)
	return int32(atoi)
}

// 将字符串转换成浮点型
// @param str 字符串
func StringToFloat64Must(str string) float64 {
	atoi, _ := strconv.ParseFloat(str, 64)
	return atoi
}

// 将字符串转换成浮点型
// @param str 字符串
func StringToFloat32Must(str string) float32 {
	atoi, _ := strconv.ParseFloat(str, 32)
	return float32(atoi)
}
