package es

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// 批量创建文档参数
// @author fuzeyu
// @date 2025/2/25
type DocCreateParams struct {
	items     []*DocCreateItemParams
	indexName string
}

// 添加一条记录
// @param id 数据id
// @param data 数据
func (obj *DocCreateParams) Add(id int64, data interface{}) error {
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	obj.items = append(obj.items, &DocCreateItemParams{
		Id:      id,
		Content: string(marshal),
	})
	return nil
}

// 返回读取结构
func (obj *DocCreateParams) Reader() (io.Reader, error) {
	if len(obj.items) == 0 {
		return nil, errors.New("没有添加有效数据")
	}
	var buffers = new(strings.Builder)
	for _, item := range obj.items {
		buffers.WriteString(`{"index": {"_index": "` + obj.indexName + `", "_id": "` + fmt.Sprintf("%d", item.Id) + `"}}` + "\n")
		buffers.WriteString(item.Content + "\n")
	}
	return strings.NewReader(buffers.String()), nil
}

func NewDocCreateParams(indexName string) *DocCreateParams {
	c := &DocCreateParams{
		indexName: indexName,
		items:     make([]*DocCreateItemParams, 0, 5),
	}
	return c
}

// 创建记录
type DocCreateItemParams struct {
	Id      int64
	Content string
}
