package es

import (
	"encoding/json"
	"errors"
	"fmt"
)

// 查询结果
type EsResponse struct {
	Took         int                        `json:"took"`
	TimedOut     bool                       `json:"timed_out"`
	Shards       Shards                     `json:"_shards"`
	Hits         HitsWrapper                `json:"hits"`
	Aggregations map[string]json.RawMessage `json:"aggregations"`
}

// 获取聚合结果
func (obj *EsResponse) Group() (*GroupResult, error) {

	var msg []byte
	for _, message := range obj.Aggregations {
		msg = message
	}
	if len(msg) < 0 {
		return nil, errors.New(fmt.Sprintf("未查询到聚合结果信息"))
	}

	var ret *GroupResult
	if err := json.Unmarshal(msg, &ret); err != nil {
		return nil, errors.New(fmt.Sprintf("解析聚合结果失败 %s", err.Error()))
	}
	return ret, nil
}

// 获取聚合值
func (obj *EsResponse) AggsValue() (float64, error) {

	var msg []byte
	for _, message := range obj.Aggregations {
		msg = message
	}
	if len(msg) < 0 {
		return 0, errors.New(fmt.Sprintf("未查询到聚合结果信息"))
	}

	var jsonValues GroupValue
	if err := json.Unmarshal(msg, &jsonValues); err != nil {
		return 0, errors.New(fmt.Sprintf("解析聚合结果失败 %s", err.Error()))
	}
	return jsonValues.Value, nil
}

type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type HitsWrapper struct {
	Total    HitsWrapperTotal `json:"total"`
	MaxScore float64          `json:"max_score"`
	Hits     []Hit            `json:"hits"`
}
type HitsWrapperTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type Hit struct {
	Index  string          `json:"_index"`
	ID     string          `json:"_id"`
	Score  float64         `json:"_score"`
	Source json.RawMessage `json:"_source"` // 使用 RawMessage 延迟解析
}

// 分词器
type ProAnalyzer string

const (
	ProAnalyzerIKMaxWord ProAnalyzer = "ik_max_word" // 提供更高的召回率，但可能增加索引大小和查询时间
	ProAnalyzerIKSmart   ProAnalyzer = "ik_smart"    // 更高效，适合对语义完整性要求较高的场景
	ProAnalyzerStand     ProAnalyzer = "standard"
)

// Token 分词结果结构体
type Token struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

// TokenList 多个分词结果的集合
type TokenList struct {
	Tokens []Token `json:"tokens"`
}

// 聚合结果
type GroupResult struct {
	DocCountErrorUpperBound int                  `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int                  `json:"sum_other_doc_count"`
	Buckets                 []GroupResultBuckets `json:"buckets"`
}

type GroupResultBuckets struct {
	Key      int `json:"key"`
	DocCount int `json:"doc_count"`
}

type GroupValue struct {
	Value float64 `json:"value"`
}
