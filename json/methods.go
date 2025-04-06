package json

import (
	"github.com/bytedance/sonic"
	"io"
)

var _default = sonic.ConfigDefault

// 获取JSON解析信息
// @param data 解析对象
func NewDecoder(data io.Reader) sonic.Decoder {
	return _default.NewDecoder(data)
}

// 获取JSON生成信息
// @param writer 编译对象
func NewEncoder(writer io.Writer) sonic.Encoder {
	return _default.NewEncoder(writer)
}

// 将对象直接编译成JSON
// @param data 编译对象
func Marshal(data any) ([]byte, error) {
	return _default.Marshal(data)
}

// 将JSON编译成对象
// @param data 编译对象
// @param v 解析对象
func Unmarshal(buffers []byte, data interface{}) error {
	return _default.Unmarshal(buffers, data)
}
