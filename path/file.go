package path

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 在一个地址中写入文件
// @param path 文件路径
// @param content 文件内容
func WriterFile(path string, content string) error {

	// 创建文件夹
	err := FolderMustCreate(filepath.Dir(path))
	if err != nil {
		return err
	}

	// 创建文件
	create, err := os.Create(path)
	if err != nil {
		return err
	}
	defer create.Close()

	// 保存内容
	reader := strings.NewReader(content)
	if _, err := io.Copy(create, reader); err != nil {
		return err
	}
	return nil
}

// 检查文件是否存在
// @param path 文件路径
func FileExist(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	if stat != nil && stat.IsDir() == false {
		return true
	}
	return false
}
