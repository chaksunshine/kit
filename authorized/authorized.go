package authorized

import (
	"fmt"
	"github.com/chaksunshine/kit/encryption"
	"github.com/chaksunshine/kit/numbers"
	"math"
	"time"
)

// 密钥生成
// @author fuzeyu
// @date 2025/4/5
type authorized struct {

	// 密钥信息
	aes *encryption.AesCry

	// 允许的请求误差数
	allowSecondLimit float64
}

// 初始化密钥
// @param securityKey 密钥
func (obj *authorized) initSecurityKey(securityKey string) error {

	aes, err := encryption.NewAes(securityKey)
	if err != nil {
		return err
	}
	obj.aes = aes

	obj.allowSecondLimit = 15

	return nil
}

// 创建一个新的请求密钥
func (obj *authorized) newRequestToken() (string, error) {
	var seconds = fmt.Sprintf("%d", time.Now().Unix())
	return obj.aes.Encrypt([]byte(seconds))
}

// 验证鉴权密钥
// @param token 鉴权密钥
func (obj *authorized) validateRequestToken(token string) error {

	// 是否解密成功
	decrypt, err := obj.aes.Decrypt(token)
	if err != nil {
		return ErrValidateAuthorizedFail
	}

	// 检查是否超时
	var unix = time.Unix(numbers.StringToInt64Must(string(decrypt)), 0)
	now := time.Now()
	if math.Abs(float64(now.Unix()-unix.Unix())) > obj.allowSecondLimit {
		return ErrValidateAuthorizedFail
	}
	return nil
}
