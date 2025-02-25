package es

import "github.com/elastic/go-elasticsearch/v8"

// 配置
type Config struct {
	Address []string

	instance *elasticsearch.Client
}

// 创建客户端
func (obj *Config) NewClient() (*elasticsearch.Client, error) {

	if obj.instance != nil {
		return obj.instance, nil
	}

	// 获取客户端
	config := elasticsearch.Config{
		Addresses: obj.Address,
	}
	client, err := elasticsearch.NewClient(config)
	if err != nil {
		return nil, err
	}
	obj.instance = client
	return client, nil
}
