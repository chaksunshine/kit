package es

import (
	"bytes"
	"encoding/json"
	"io"
)

// 查询数
// @author fuzeyu
// @date 2025/2/26
type Query struct {
	docName string

	conditions map[string]interface{}
}

// 输出查询条件
func (obj *Query) String() string {
	marshal, _ := json.Marshal(obj.conditions)
	return string(marshal)
}

// 获取查询内容
func (obj *Query) Reader() (io.Reader, error) {
	marshal, err := json.Marshal(obj.conditions)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(marshal), nil
}

// 初始化一个字段
// @param container 容器
// @param fields 字段信息
func (obj *Query) initQueryFields(container map[string]interface{}, fields string, kindMap queryKind) {
	if _, ok := container[fields]; ok == false {
		switch kindMap {
		case queryKindArray:
			container[fields] = make([]interface{}, 0)

		case queryKindMap:
			container[fields] = make(map[string]interface{})
		}
	}
}

// 按照and方法初始化基础对象
// @param call 需要条件的回调条件
func (obj *Query) baseAnd(call map[string]interface{}) {

	// 构造基础条件
	obj.initQueryFields(obj.conditions, `query`, queryKindMap)
	container := obj.conditions[`query`].(map[string]interface{})

	obj.initQueryFields(container, `bool`, queryKindMap)
	container = container[`bool`].(map[string]interface{})

	obj.initQueryFields(container, `must`, queryKindArray)
	mustContainer := container[`must`].([]interface{})

	mustContainer = append(mustContainer, call)
	container[`must`] = mustContainer
}

// 按照模糊的方式查询
// @param fields 查询字段
// @param value 查询内容
func (obj *Query) Like(fields string, value string) *Query {
	obj.baseAnd(map[string]interface{}{
		"match": map[string]interface{}{
			fields: value,
		},
	})
	return obj
}

// 按照模糊的方式查询
// @param fields 查询字段
// @param value 查询内容
func (obj *Query) Equals(fields string, value string) *Query {
	obj.baseAnd(map[string]interface{}{
		"term": map[string]interface{}{
			fields: value,
		},
	})
	return obj
}

// 覆盖条件
// @param conditions 查询条件
func (obj *Query) Full(conditions map[string]interface{}) *Query {
	obj.conditions = conditions
	return obj
}

// 查询的分页信息
// @param page 页码
// @param size 每页数量
func (obj *Query) Page(page int, size int) *Query {
	if page <= 0 {
		page = 1
	}
	obj.conditions["from"] = size * (page - 1)
	obj.conditions["size"] = size
	return obj
}

// 覆盖之前的条件
// @param conditions 方式
func (obj *Query) Conditions(conditions map[string]interface{}) *Query {
	obj.conditions = conditions
	return obj
}

// @param docName 查询内容
func NewQuery(docName string) *Query {
	c := &Query{
		docName:    docName,
		conditions: make(map[string]interface{}),
	}
	return c
}

// 查询的时候,初始化的数据类型
type queryKind int

const (
	queryKindMap queryKind = iota
	queryKindArray
)
