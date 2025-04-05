package forwarding

import (
	"fmt"
	"github.com/chaksunshine/kit/catLog"
	"io"
	"net"
	"strings"
	"time"
)

// 端口转发
// @author fuzeyu
// @date 2025/4/5
type PortForwarding struct {
	config *Config
	target string

	allow map[string]bool
}

// 解析允许的IP地址
func (obj *PortForwarding) parserAllowIp() error {

	obj.allow = make(map[string]bool)
	for _, item := range obj.config.AllowIp {
		obj.allow[item] = true
	}

	return nil
}

// 处理消息端口信息
// @param accept 消息端口
func (obj *PortForwarding) handle(accept net.Conn) {

	// 初始化
	address := accept.RemoteAddr().String()
	catLog.Info(fmt.Sprintf("端口信息转发 开始 %s  =>  %s", address, obj.target))
	defer func() {
		time.Sleep(time.Millisecond * 100)
		catLog.Info(fmt.Sprintf("端口信息转发 结束 %s  =>  %s", address, obj.target))
		_ = accept.Close()
	}()
	var clientAddress = address[:strings.Index(address, ":")]

	// 验证IP地址
	if obj.allow[clientAddress] == false {
		return
	}

	// 消息转发
	dial, err := net.Dial("tcp", obj.target)
	if err != nil {
		catLog.Error(fmt.Sprintf("消息转发失败 连接到目标失败 %s %s", obj.target, err.Error()))
		return
	}
	defer dial.Close()

	// 双向转发
	go func() {
		_, _ = io.Copy(dial, accept)
	}()
	_, _ = io.Copy(accept, dial)
}

// 开始转发
func (obj *PortForwarding) Start() error {

	// 解析允许的IP地址
	if err := obj.parserAllowIp(); err != nil {
		return err
	}

	// 创建端口
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", obj.config.ListenPort))
	if err != nil {
		return err
	}
	defer listener.Close()

	// 监听
	for {
		accept, err := listener.Accept()
		if err != nil {
			catLog.Error(fmt.Sprintf("端口转发中获取信息失败 %s", err.Error()))
			continue
		}
		go obj.handle(accept)
	}
}

// 端口转发
// @param config *Config
func NewPortForwarding(config *Config) *PortForwarding {
	c := &PortForwarding{
		config: config,
		allow:  make(map[string]bool),
		target: fmt.Sprintf("%s:%d", config.TargetHost, config.TargetPort),
	}
	return c
}
