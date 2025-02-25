package id

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
)

// 雪花唯一serial
// @author fuzeyu
// @date 2025/2/25
type serial struct {
	node  *snowflake.Node
	cache chan int64
}

// 获取下一个
func (obj *serial) Next() int64 {
	return <-obj.cache
}

// 循环创建
func (obj *serial) loopMarket() {
	for {
		obj.cache <- obj.node.Generate().Int64()
	}
}

// 初始化
func (obj *serial) init() {

	node, err := snowflake.NewNode(1002)
	if err != nil {
		panic(errors.New(fmt.Sprintf("初始化雪花id失败 %s", err.Error())))
	}
	obj.node = node
	go obj.loopMarket()
}

func newSerial() *serial {
	c := &serial{
		cache: make(chan int64, 15),
	}
	c.init()
	return c
}

var Serial = newSerial()
