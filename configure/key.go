package configure

import (
	"github.com/alecthomas/kingpin"
	"github.com/chaksunshine/kit/numbers"
	"os"
	"strings"
)

// 获取一个Key
// @author fuzeyu
// @date 2025/3/9
type Key struct {

	// 地址栏参数
	set map[string]*keyItem

	// 环境变量
	env map[string]string
}

// 解析
func (obj *Key) Parser() {
	for _, item := range obj.set {
		target := kingpin.Flag(item.name, item.contents).String()
		item.value = target
	}
	kingpin.Parse()
}

// 按照一个数字类型的方式获取参数
// @param key 参数信息
func (obj *Key) Int32(key string) int32 {
	return numbers.StringToInt32Must(obj.String(key))
}

// 获取一个参数信息
// @param key 参数名
func (obj *Key) String(key string) string {

	// 检查环境变量中是否存在
	val, ok := obj.env[key]
	if ok && val != "" {
		return val
	}

	// 检查地址栏参数是否会存在
	item, ok := obj.set[key]
	if !ok {
		return ""
	}
	ret := *item.value
	if ret != "" {
		return ret
	}

	return item.defaultValue
}

// 解析地址栏参数
func (obj *Key) parserEnv() {
	for _, item := range os.Environ() {
		datum := strings.Split(item, "=")
		if len(datum) == 2 && strings.HasPrefix(datum[0], EnvPrefix) {
			prefix := strings.TrimPrefix(datum[0], EnvPrefix)
			all := strings.ReplaceAll(strings.ToLower(prefix), "_", "-")
			obj.env[all] = datum[1]
		}
	}
}

// 创建一个配置信息
// @param name 名称
// @param content 内容
// @param defaultValue 默认值
func (obj *Key) Add(name string, content string, defaultValue string) {
	obj.set[name] = &keyItem{
		name:         name,
		contents:     content,
		defaultValue: defaultValue,
	}
}

func NewKey() *Key {
	c := &Key{
		set: make(map[string]*keyItem),
		env: make(map[string]string),
	}
	c.parserEnv()
	return c
}

type keyItem struct {
	value        *string
	name         string
	defaultValue string
	contents     string
}

// 按照服务发现模式注册常用信息
func RegisterServiceDiscovery() *Key {
	key := NewKey()
	key.Add("master-address", "主控服务地址 示例: server:32000", "127.0.0.1:32000")
	key.Add("service-name", "当前项目服务名称，用于获取服务配置", "users")
	key.Add("listen-port", "当前项目需要监听的服务信息", "http.32001,grpc.32002")
	key.Parser()
	return key
}

// 按照环境变量的方式获取
var DefaultKey = NewKey()

const (

	// 获取服务的环境变量
	EnvPrefix = `SERVICE_ENV_`
)

// 允许的协议
type Protoc string

const (
	ProtocGrpc Protoc = "grpc"
	ProtocHttp Protoc = "http"
)
