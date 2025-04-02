package etcd

// 获取节点信息
type Member struct {
	Id       uint64
	Name     string
	Address  string
	IsLeader bool
}

type Members []*Member

func (obj Members) Len() int {
	return len(obj)
}

func (obj Members) Less(i, j int) bool {
	return obj[i].IsLeader && !obj[j].IsLeader
}

func (obj Members) Swap(i, j int) {
	obj[i], obj[j] = obj[j], obj[i]
}
