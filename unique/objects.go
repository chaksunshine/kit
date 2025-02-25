package unique

import (
	"bytes"
	"encoding/gob"
)

// 深度拷贝
// @param src 原始数据
func DeepCopy[T any](src T) (T, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(src); err != nil {
		return src, err
	}

	var dst T
	dec := gob.NewDecoder(&buf)
	if err := dec.Decode(&dst); err != nil {
		return src, err
	}
	return dst, nil
}
