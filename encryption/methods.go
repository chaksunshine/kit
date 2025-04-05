package encryption

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"net/url"
)

// 计算md5
// @param content 内容
func Md5(content string) string {
	hash := md5.New()
	hash.Write([]byte(content))
	sum := hash.Sum(nil)
	return hex.EncodeToString(sum)
}

// 计算Sha1
// @param content 内容
func Sha1(content string) string {
	hash := sha1.New()
	hash.Write([]byte(content))
	sum := hash.Sum(nil)
	return hex.EncodeToString(sum)
}

// 计算Url编码
// @param content 内容
func UrlEncode(content string) string {
	return url.QueryEscape(content)
}

// 计算Url解码
// @param content 内容
func UrlDecode(content string) string {
	unescape, err := url.QueryUnescape(content)
	if err != nil {
		return ""
	}
	return unescape
}

// 按照base64编码
// @param content 内容
func Base64Encode(content string) string {
	return base64.StdEncoding.EncodeToString([]byte(content))
}

// 按照base64解码
// @param content 内容
func Base64Decode(content string) string {
	decodeString, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return ""
	}
	return string(decodeString)
}
