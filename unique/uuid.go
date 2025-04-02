package unique

import (
	"crypto/rand"
	"fmt"
)

// 实现一个随机生成uuid的方法
func NewUUID() string {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return ""
	}
	randomBytes[6] = (randomBytes[6] & 0x0f) | 0x40
	randomBytes[8] = (randomBytes[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		randomBytes[0:4], randomBytes[4:6], randomBytes[6:8],
		randomBytes[8:10], randomBytes[10:])
}
