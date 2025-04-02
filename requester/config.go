package requester

// 请求配置
type Config struct {

	// 是否要格式化header头
	FormatHeader bool
}

var defaultConfig = &Config{
	FormatHeader: true,
}
