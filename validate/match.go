package validate

import (
	"errors"
	"fmt"
	"strconv"
)

// 验证
// @author fuzeyu
// @date 2025/2/20
type match struct {
}

// 检查是否是IPv4地址
// @param host 节点信息
func (obj *match) IsIpV4(host string) bool {
	return RegexpIPV4.MatchString(host)
}

// 检查字符串
// @param name 字段名
// @param value 字段值
// @param max 最大长度
// @param min 最小长度
func (obj *match) IsString(name, value string, min, max int) error {
	length := len(value)
	if length > max {
		return errors.New(fmt.Sprintf("%s长度不能超过%d个字符", name, max))
	}
	if length < min {
		return errors.New(fmt.Sprintf("%s长度不能少于%d个字符", name, min))
	}
	return nil
}

// 检查是否为字符串数字
// @param name 字段名
// @param value 字段值
// @param max 最大值
// @param min 最小值
// @return error
func (obj *match) IsStringNumber(name, value string) error {
	if RegexpStringNumber.MatchString(value) == false {
		return errors.New(fmt.Sprintf("%s不是数字格式", name))
	}
	return nil
}

// 检查整数
// @param name 字段名
// @param value 字段值
// @param max 最大值
// @param min 最小值
// @return error
func (obj *match) IsInt(name string, value int, min, max int) error {
	if value > max {
		return errors.New(fmt.Sprintf("%s不能大于%d", name, max))
	}
	if value < min {
		return errors.New(fmt.Sprintf("%s不能小于%d", name, min))
	}
	return nil
}

// 检查整数
// @param name 字段名
// @param value 字段值
// @param max 最大值
// @param min 最小值
// @return error
func (obj *match) IsInt64(name string, value int64, min, max int64) error {
	if value > max {
		return errors.New(fmt.Sprintf("%s不能大于%d", name, max))
	}
	if value < min {
		return errors.New(fmt.Sprintf("%s不能小于%d", name, min))
	}
	return nil
}

// 检查整数
// @param name 字段名
// @param value 字段值
// @param max 最大值
// @param min 最小值
// @return error
func (obj *match) IsInt32(name string, value int32, min, max int32) error {
	if value > max {
		return errors.New(fmt.Sprintf("%s不能大于%d", name, max))
	}
	if value < min {
		return errors.New(fmt.Sprintf("%s不能小于%d", name, min))
	}
	return nil
}

// 检查大陆身份证号
// @param idCard 身份证号码
func (obj *match) IsZhIDCard(idCard string) bool {
	if !RegexpChinaIDCard.MatchString(idCard) {
		return false
	}
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checkCodes := "10X98765432"

	sum := 0
	for i := 0; i < 17; i++ {
		num, _ := strconv.Atoi(string(idCard[i]))
		sum += num * weights[i]
	}
	expectedCheckCode := checkCodes[sum%11]
	actualCheckCode := idCard[17]
	return expectedCheckCode == actualCheckCode
}

// 检查是否是中国大陆手机号
// @param phoneNumber 手机号码
func (obj *match) IsZhCnPhone(phoneNumber string) bool {
	return RegexpChinaPhone.MatchString(phoneNumber)
}

// 检查是否是邮箱
// @param email 邮箱
func (obj *match) IsEmailAddress(email string) bool {
	return RegexpEmailAddress.MatchString(email)
}

var Match = new(match)
