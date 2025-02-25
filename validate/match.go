package validate

import "regexp"

// 验证
// @author fuzeyu
// @date 2025/2/20
type match struct {
}

// 检查是否是IPv4地址
// @param host 节点信息
func (obj *match) IsIpV4(host string) bool {
	return RegexpMatchIPv4.MatchString(host)
}

var MatchInstance = new(match)

// 正则匹配IP地址
var RegexpMatchIPv4 = regexp.MustCompile("^((25[0-5]|2[0-4]\\d|[01]?\\d\\d?)\\.){3}(25[0-5]|2[0-4]\\d|[01]?\\d\\d?)$")
