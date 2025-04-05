package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// Aes加解密
type AesCry struct {
	aesKey []byte
}

// 验证密钥信息是否正确
// @param key Key信息
func (obj *AesCry) parserHashedKey(key string) error {
	lengths := 60
	if len(key) < lengths {
		return errors.New(fmt.Sprintf("密钥长度不能少于 %d 位", lengths))
	}
	sum256 := sha256.Sum256([]byte(key))
	obj.aesKey = sum256[:32]
	return nil
}

// AES加密核心函数
// @param  plaintext []byte 待加密的明文
// @param key string 密钥
func (obj *AesCry) Encrypt(value []byte) (string, error) {

	// 生成随机IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// 创建加密块
	block, err := aes.NewCipher(obj.aesKey)
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
func (obj *AesCry) Decrypt(encodedText string) ([]byte, error) {

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
	block, err := aes.NewCipher(obj.aesKey)
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
func (obj *AesCry) pKCS7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	return append(src, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

// PKCS7去填充
func (obj *AesCry) pKCS7Unpadding(src []byte) ([]byte, error) {
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

// @param serverKey 加密密钥
func NewAes(key string) (*AesCry, error) {
	data := &AesCry{}
	if err := data.parserHashedKey(key); err != nil {
		return nil, err
	}
	return data, nil
}
