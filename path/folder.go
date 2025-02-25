package path

import (
	"errors"
	"fmt"
	"os"
)

// 目录是否存在
// @param dir 文件夹
func FolderMustCreate(dir string) error {

	// 检查目录是否存在,如果不存在就必须创建
	stat, _ := os.Stat(dir)
	if stat == nil {
		return os.MkdirAll(dir, os.ModePerm)
	}

	// 检查是否是目录
	if stat.IsDir() {
		return nil
	}
	return errors.New(fmt.Sprintf("%s 存在,但不是目录", dir))
}
