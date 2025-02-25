package es

import (
	"errors"
	"fmt"
	"github.com/chaksunshine/kit/thread"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"io"
	"net/http"
	"strings"
	"time"
)

// 文档管理
type Doc struct {
	config *Config
	client *elasticsearch.Client
}

// 执行
func (obj *Doc) connect() error {

	// 创建客户端
	if client, err := obj.config.NewClient(); err != nil {
		return err
	} else {
		obj.client = client
	}
	return nil
}

// 批量创建
// @param docName 文档名称
func (obj *Doc) Create(docName string, createParams *DocCreateParams) error {

	// 保存创建参数
	reader, err := createParams.Reader()
	if err != nil {
		return err
	}

	// 发送批量请求
	req := esapi.BulkRequest{
		Body: reader,
	}
	buffers, err := req.Do(thread.CtxRequest(), obj.client)
	if err != nil {
		return err
	}
	defer buffers.Body.Close()
	if buffers.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("批量创建失败 %s %d\n %s", docName, buffers.StatusCode, buffers.String()))
	}

	// 检查是否存在错误
	all, err := io.ReadAll(buffers.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("批量创建失败 %s %d\n %s", docName, buffers.StatusCode, buffers.String()))
	}
	if text := string(all); strings.HasPrefix(text, "{\"errors\":true,") {
		return errors.New(fmt.Sprintf("批量创建失败 %s \n %s", docName, text))
	}
	return nil
}

// 删除文档
// @param name 名称
// @param id 要删除的id
func (obj *Doc) Drop(name string, index ...int64) error {

	if len(index) == 0 {
		return errors.New("没有添加有效数据")
	}
	var buffers = new(strings.Builder)
	for _, item := range index {
		buffers.WriteString(fmt.Sprintf(`{"delete": {"_index": "%s", "_id": %d}}`, name, item) + "\n")
	}

	// 执行 Bulk API 请求
	msg := buffers.String()
	req := esapi.BulkRequest{
		Body:    strings.NewReader(msg),
		Pretty:  true,                            // 格式化输出
		Refresh: "true",                          // 刷新索引以立即可见
		Timeout: time.Duration(30) * time.Second, // 超时时间
	}
	do, err := req.Do(thread.CtxRequest(), obj.client)
	if err != nil {
		return err
	}
	defer do.Body.Close()

	if do.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("批量删除失败 %s %d\n %s", name, do.StatusCode, do.String()))
	}
	return nil
}

// @param config 参数信息
func NewDoc(config *Config) (*Doc, error) {
	doc := &Doc{
		config: config,
	}
	return doc, doc.connect()
}
