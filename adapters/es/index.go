package es

import (
	"errors"
	"fmt"
	"github.com/chaksunshine/kit/thread"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"net/http"
)

// es Index 操作助手
// @author fuzeyu
// @date 2025/2/25
type Index struct {
	cfg *Config

	client *elasticsearch.Client
}

// 执行
func (obj *Index) connect() error {

	// 创建客户端
	if client, err := obj.cfg.NewClient(); err != nil {
		return err
	} else {
		obj.client = client
	}
	return nil
}

// 创建索引
// @param names 索引名称
// @param params 索引参数
func (obj *Index) Create(names string, params *IndexCreateParams) error {

	// 检查是否存在
	exists, err := obj.client.Indices.Exists([]string{names})
	if err != nil {
		return errors.New(fmt.Sprintf("创建索引失败 %s %s", names, err.Error()))
	}

	// 不存在
	if exists.StatusCode != http.StatusNotFound {
		return nil
	}

	// 创建
	request := esapi.IndicesCreateRequest{
		Index: names,
		Body:  params.Reader(),
	}
	buffers, err := request.Do(thread.CtxRequest(), obj.client)
	if err != nil {
		return err
	}
	defer buffers.Body.Close()

	if buffers.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("创建索引失败 %s %d \n%s", names, buffers.StatusCode, buffers.String()))
	}
	return nil
}

// 删除索引
// @param name
func (obj *Index) Drop(name string) error {

	deleteReq := esapi.IndicesDeleteRequest{
		Index: []string{name},
	}
	resp, err := deleteReq.Do(thread.CtxRequest(), obj.client)
	if err != nil {
		return errors.New(fmt.Sprintf("删除索引失败 %s %s", name, err.Error()))
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("删除索引失败 %s %d", name, resp.StatusCode))
	}
	return nil
}

func NewIndex(cfg *Config) (*Index, error) {
	index := &Index{
		cfg: cfg,
	}
	return index, index.connect()
}
