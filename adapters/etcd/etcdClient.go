package etcd

import (
	"errors"
	"fmt"
	"github.com/chaksunshine/kit/catLog"
	"github.com/chaksunshine/kit/configure"
	"github.com/chaksunshine/kit/thread"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sort"
	"strings"
	"time"
)

// etcd客户端
// @author fuzeyu
// @date 2025/3/6
type EtcdClient struct {
	config *Config
	client *clientv3.Client
}

// 获取客户端
func (obj *EtcdClient) Client() *clientv3.Client {
	return obj.client
}

// 执行
func (obj *EtcdClient) connect() error {

	var nodes = make([]string, 0, 5)
	for _, item := range strings.Split(obj.config.Nodes, ",") {
		if len(item) > 0 {
			nodes = append(nodes, item)
		}
	}
	if len(nodes) == 0 {
		return errors.New("请配置etcd服务器链接地址")
	}

	// 连接到etcd集群
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   nodes, // etcd集群节点地址
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		return err
	}
	obj.client = cli

	// 检查是否链接成功
	status, err := cli.Status(thread.CtxRequest(), nodes[0])
	if err != nil {
		return errors.New(fmt.Sprintf("链接到etcd集群失败 %s", err.Error()))
	}
	catLog.Info(fmt.Sprintf("链接到etcd服务成功 %s 版本: %s", obj.config.Nodes, status.Version))

	return err
}

// 获取一个key的最新值
// @param prefix 前缀信息
func (obj *EtcdClient) GetKey(prefix string) (*mvccpb.KeyValue, error) {
	get, err := obj.client.Get(thread.CtxRequest(), prefix)
	if err != nil {
		return nil, err
	}
	if get != nil && len(get.Kvs) > 0 {
		return get.Kvs[0], nil
	}
	return nil, nil
}

// 获取一个key的最新值
// @param prefix 前缀信息
func (obj *EtcdClient) GetPrefix(prefix string, size ...int) ([]*mvccpb.KeyValue, error) {

	var findSize = 2000
	if len(size) > 0 {
		findSize = size[0]
	}

	get, err := obj.client.Get(thread.CtxRequest(), prefix, clientv3.WithPrefix(), clientv3.WithKeysOnly(), clientv3.WithLimit(int64(findSize)))
	if err != nil {
		return nil, err
	}
	if get != nil && len(get.Kvs) > 0 {
		return get.Kvs, nil
	}
	return nil, nil
}

// 获取节点信息
func (obj *EtcdClient) Members() (Members, error) {

	// 获取客户端信息
	request := thread.CtxRequest()
	resp, err := obj.client.MemberList(request)
	if err != nil {
		return nil, err
	}

	var ret = make(Members, 0, len(resp.Members))
	for _, item := range resp.Members {
		ret = append(ret, &Member{
			Id:      item.ID,
			Name:    item.Name,
			Address: item.ClientURLs[0],
		})
	}
	if len(ret) <= 0 {
		return nil, errors.New("没有获取到存活的etcd节点")
	}

	// 查询leader节点信息
	status, err := obj.client.Status(request, ret[0].Address)
	if err != nil {
		return nil, errors.New("无法获取到leader节点信息")
	}
	for _, member := range ret {
		member.IsLeader = member.Id == status.Leader
	}
	sort.Sort(ret)
	return ret, nil
}

// 关闭客户端
func (obj *EtcdClient) Close() {
	obj.client.Close()
}

func NewEtcdClient(config *Config) (*EtcdClient, error) {
	c := &EtcdClient{
		config: config,
	}
	return c, c.connect()
}

// 按照字符串的形式创建一个etcd客户端
// @param str etcd服务器链接配置
func NewsEtcdClientByString(str string) (*EtcdClient, error) {

	var etcdConfig *Config
	buffers := []byte(str)
	if err := configure.LoadBufferConfigure(buffers, &etcdConfig); err != nil {
		return nil, err
	}
	return NewEtcdClient(etcdConfig)
}
