package configure

import (
	"errors"
	"fmt"
	"github.com/chaksunshine/kit/define"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

// 加载配置文件信息
// @param name 名称
// @param data 对象信息
func LoadFileConfigure(name string, data interface{}) error {

	// 加载配置
	buffers, paths := ReadFile(name)
	if len(buffers) == 0 {
		return errors.New(fmt.Sprintf("没有读取到配置文件，请检查\n%s", strings.Join(paths, "\n")))
	}

	// 解析
	return LoadBufferConfigure(buffers, data)
}

// 按照Buffer的方式解析配置文件
func LoadBufferConfigure(buffers []byte, data interface{}) error {
	if err := yaml.Unmarshal(buffers, data); err != nil {
		return errors.New(fmt.Sprintf("解析配置文件失败 %s", err.Error()))
	}
	return nil
}

// 读取一次文件
// @param name 文件信息
func ReadFile(name string) ([]byte, []string) {
	var paths = []string{
		filepath.Join(define.RootPath, fmt.Sprintf("%s.yml", name)),
		filepath.Join(define.RootPath, fmt.Sprintf("%s.yaml", name)),
		filepath.Join(define.RootPath, "config", fmt.Sprintf("%s.yml", name)),
		filepath.Join(define.RootPath, "config", fmt.Sprintf("%s.yaml", name)),
	}
	for _, item := range paths {
		buffers, err := os.ReadFile(item)
		if err == nil && len(buffers) > 0 {
			return buffers, nil
		}
	}
	return nil, paths
}
