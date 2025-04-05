package forwarding

// 端口转发
type Config struct {

	// 监听的端口
	ListenPort int32 `yaml:"listenPort"`

	// 转发的目标
	TargetHost string `yaml:"targetHost"`
	TargetPort int32  `yaml:"targetPort"`

	// 授权访问的IP地址
	AllowIp []string `yaml:"allowIp"`
}
