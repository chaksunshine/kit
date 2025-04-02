package define

import "os"

// 获取项目目录
var RootPath string

// 资源目录
var DataPath string

func init() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	RootPath = dir

	DataPath = RootPath + "/data"
}
