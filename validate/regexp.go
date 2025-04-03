package validate

import "regexp"

// 检查是否是IPV4的IP地址
var RegexpIPV4 = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

// 中国大陆手机号
var RegexpChinaPhone = regexp.MustCompile(`^1[3-9]\d{9}$`)

// 中国大陆身份证格式
var RegexpChinaIDCard = regexp.MustCompile(`^\d{17}[\dXx]$`)

// 邮箱格式
var RegexpEmailAddress = regexp.MustCompile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`)

// 数字格式
var RegexpStringNumber = regexp.MustCompile(`^[0-9]+$`)
