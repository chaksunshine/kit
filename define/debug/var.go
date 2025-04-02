package debug

import (
	"os"
	"runtime"
)

// 当前是否是windows平台
var IsWindows = runtime.GOOS == `windows`

// 检查是否是在idea中运行环境
var IsIdea = os.Getenv(`JB_GOLAND`) == `dFX06CJErVG2DOL179nAHzmQTNxbtkWy5ulRSUesgafwi3MZPI4Kqcjv8ohpYBlr`
