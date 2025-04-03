package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

type aesCry struct {
}

// AES加密核心函数
// @param  plaintext []byte 待加密的明文
// @param key string 密钥
func (obj *aesCry) Encrypt(value []byte, key string) (string, error) {

	if len(key) < 32 {
		return "", errors.New("密钥长度过短")
	}

	hashedKey := sha256.Sum256([]byte(key))
	aesKey := hashedKey[:32] // 使用SHA256生成32字节密钥

	// 生成随机IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// 创建加密块
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	// PKCS7填充
	value = obj.pKCS7Padding(value, aes.BlockSize)
	ciphertext := make([]byte, len(value))

	// CBC模式加密
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, value)

	// 合并IV和密文
	combined := append(iv, ciphertext...)
	return base64.StdEncoding.EncodeToString(combined), nil
}

// AES解密核心函数
// @param  encodedText []byte 编码密钥
// @param key string 密钥
func (obj *aesCry) Decrypt(encodedText string, key string) ([]byte, error) {
	// 处理密钥
	hashedKey := sha256.Sum256([]byte(key))
	aesKey := hashedKey[:32]

	// 解码Base64
	combined, err := base64.StdEncoding.DecodeString(encodedText)
	if err != nil {
		return nil, err
	}

	// 分离IV和密文
	if len(combined) < aes.BlockSize {
		return nil, errors.New("invalid ciphertext")
	}
	iv := combined[:aes.BlockSize]
	ciphertext := combined[aes.BlockSize:]

	// 创建解密块
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	// CBC模式解密
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// 去除PKCS7填充
	return obj.pKCS7Unpadding(plaintext)
}

// PKCS7填充
func (obj *aesCry) pKCS7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	return append(src, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

// PKCS7去填充
func (obj *aesCry) pKCS7Unpadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, errors.New("invalid padding")
	}
	un := int(src[length-1])
	if un > length {
		return nil, errors.New("invalid padding")
	}
	return src[:length-un], nil
}

var AesCry = new(aesCry)
