package redis

import (
	"fmt"
	"strings"
)

// 节点信息
type Node struct {
	ID         string // 节点 ID
	Host       string // 节点主机地址
	Port       int    // 节点端口
	IsMaster   bool   // 是否为主节点
	MasterHost string // 父节点主机地址（仅适用于从节点）
	MasterPort int    // 父节点端口（仅适用于从节点）
}

type NodeSet []Node

func (obj NodeSet) Len() int {
	return len(obj)
}

func (obj NodeSet) Less(i, j int) bool {
	if obj[i].IsMaster && !obj[j].IsMaster {
		return true
	} else if !obj[i].IsMaster && obj[j].IsMaster {
		return false
	}
	return obj[i].Port < obj[j].Port
}

func (obj NodeSet) Swap(i, j int) {
	obj[i], obj[j] = obj[j], obj[i]
}

func (obj NodeSet) String() string {
	var bf = new(strings.Builder)
	for _, item := range obj {
		if item.IsMaster {
			bf.WriteString(fmt.Sprintf("%s:%d\t\tmaster \n", item.Host, item.Port))
		} else {
			bf.WriteString(fmt.Sprintf("%s:%d\t\tslave -> %s:%d\t\tmaster \n", item.Host, item.Port, item.MasterHost, item.MasterPort))
		}
	}
	return bf.String()
}
