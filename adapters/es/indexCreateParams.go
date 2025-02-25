package es

import (
	"bytes"
	"encoding/json"
	"io"
)

// 创建index参数信息
type IndexCreateMapperProParams struct {
	Type     ProType `json:"type"`
	Analyzer string  `json:"analyzer,omitempty"`
}

// 字段类型
type ProType string

const (
	ProTypeText    ProType = "text"    // 常规字段,按照空格分割建立索引
	ProTypeKeyword ProType = "keyword" // 原始数据,不会作为拆分
	ProTypeDate    ProType = "date"    // 日期类型
	ProTypeInteger ProType = "integer" // 类型
	ProTypeFloat   ProType = "float"   // 浮点型
	ProTypeBoolean ProType = "boolean" // 布尔类型
)

// 分词器
type ProAnalyzer string

const (
	ProAnalyzerIKMaxWord ProAnalyzer = "ik_max_word" // 提供更高的召回率，但可能增加索引大小和查询时间
	ProAnalyzerIKSmart   ProAnalyzer = "ik_smart"    // 更高效，适合对语义完整性要求较高的场景
	ProAnalyzerStand     ProAnalyzer = "standard"
)

//text、keyword、date、integer、float、boolean

// 创建index参数信息
type IndexCreateParams struct {

	// 分区设置
	Settings IndexCreateParamsSettings `json:"settings"`

	// 字段信息
	Mappings IndexCreateParamsMapping `json:"mappings"`
}
type IndexCreateParamsSettings struct {
	NumberOfShards   int `json:"number_of_shards"`
	NumberOfReplicas int `json:"number_of_replicas"`
}
type IndexCreateParamsMapping struct {
	Properties map[string]*IndexCreateMapperProParams `json:"properties"`
}

// 创建字段
// @param name 字段名称
// @param types 字段类型
// @param analyzer 分词器
func (obj *IndexCreateParams) AddProperties(name string, types ProType, analyzer ...ProAnalyzer) *IndexCreateParams {

	// 执行默认分词器
	var analyzerValue string
	if types == ProTypeText {
		analyzerValue = string(ProAnalyzerStand)

		if len(analyzer) > 0 {
			analyzerValue = string(analyzer[0])
		}
	}

	// 创建参数
	obj.Mappings.Properties[name] = &IndexCreateMapperProParams{
		Type:     types,
		Analyzer: analyzerValue,
	}
	return obj
}

// 获取io对象
func (obj *IndexCreateParams) Reader() io.Reader {
	data, _ := json.Marshal(obj)
	return bytes.NewReader(data)
}

// 格式化索引参数
// @param shared 分片
// @param replicas 副本
// @param fields 字段信息
func NewIndexCreateParams(shards, replicas int) *IndexCreateParams {
	return &IndexCreateParams{
		Settings: IndexCreateParamsSettings{
			NumberOfShards:   shards,
			NumberOfReplicas: replicas,
		},
		Mappings: IndexCreateParamsMapping{
			Properties: make(map[string]*IndexCreateMapperProParams),
		},
	}
}
