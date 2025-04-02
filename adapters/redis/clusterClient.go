package redis

import (
	"github.com/chaksunshine/kit/thread"
	"github.com/redis/go-redis/v9"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 获取集群配置
// @author fuzeyu
// @date 2025/2/18
type ClusterClient struct {
	*Client

	cfg *Config
}

// 链接到集群中
func (obj *ClusterClient) connect() (*Client, error) {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:           obj.cfg.Node,
		Password:        obj.cfg.Pwd,
		PoolSize:        obj.cfg.PoolSize,
		ConnMaxLifetime: time.Second * time.Duration(obj.cfg.ConnectTimeout),
	})
	if err := rdb.Ping(thread.CtxRequestLocal()).Err(); err != nil {
		return nil, err
	}
	return newClient(rdb), nil
}

// 获取集群节点信息
func (obj *ClusterClient) ClusterNode() (NodeSet, error) {

	result, err := obj.ClusterNodes(thread.CtxRequestLocal()).Result()
	if err != nil {
		return nil, err
	}

	row := strings.Split(result, "\n")
	var nodes = make(NodeSet, 0, len(row)-1)
	for _, line := range row {
		if line == "" {
			continue
		}

		// 按空格分割字段
		fields := strings.Fields(line)
		if len(fields) < 8 {
			continue
		}

		// 提取节点基本信息
		id := fields[0]       // 节点 ID
		addr := fields[1]     // 主机地址和端口
		role := fields[2]     // 节点角色 (master/slave)
		masterID := fields[3] // 主节点 ID（如果是从节点）

		// 解析主机地址和端口
		hostPort := strings.Split(addr, "@")[0]
		hostPortParts := strings.Split(hostPort, ":")
		if len(hostPortParts) != 2 {
			continue
		}
		host := hostPortParts[0]
		port, _ := strconv.Atoi(hostPortParts[1])
		isMaster := strings.HasSuffix(role, "master")

		// 如果是从节点，提取主节点的 IP 地址
		var masterHost string
		var masterPort int
		if !isMaster {

			// TODO 性能优化
			for _, otherLine := range row {
				otherFields := strings.Fields(otherLine)
				if len(otherFields) >= 2 && otherFields[0] == masterID {
					masterAddr := otherFields[1]
					masterHostPort := strings.Split(masterAddr, "@")[0]
					masterHostPortParts := strings.Split(masterHostPort, ":")
					if len(masterHostPortParts) == 2 {
						masterHost = masterHostPortParts[0]
						masterPort, _ = strconv.Atoi(masterHostPortParts[1])
					}
					break
				}
			}
		}

		// 构造节点对象
		node := Node{
			ID:         id,
			Host:       host,
			Port:       port,
			IsMaster:   isMaster,
			MasterHost: masterHost,
			MasterPort: masterPort,
		}
		nodes = append(nodes, node)
	}

	sort.Sort(nodes)
	return nodes, nil
}

// @param cfg 配置信息
func NewClusterClient(cfg *Config) (*ClusterClient, error) {

	if err := cfg.init(); err != nil {
		return nil, err
	}

	c := &ClusterClient{
		cfg: cfg,
	}
	connect, err := c.connect()
	if err != nil {
		return nil, err
	}
	c.Client = connect
	return c, nil
}
