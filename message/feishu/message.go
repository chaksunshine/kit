package feishu

import (
	"encoding/json"
	"fmt"
	"strings"
)

// 创建消息
// @author fuzeyu
// @date 2025/4/2
type Message struct {
	title   string
	label   [][]string
	content string
}

// 设置标题
// @param title 标题
func (obj *Message) Title(title string) *Message {
	obj.title = title
	return obj
}

// 设置标签
// @param key 标签key
func (obj *Message) Label(key string, value string) *Message {
	obj.label = append(obj.label, []string{key, value})
	return obj
}

// 设置正文
// @param key 标签key
func (obj *Message) Content(content string) *Message {
	obj.content = content
	return obj
}

// 解析消息
func (obj *Message) parser() string {

	var msg = new(strings.Builder)

	if obj.title != "" {
		msg.WriteString(fmt.Sprintf("【%s】\n", obj.title))
	}

	if obj.content != "" {
		msg.WriteString(fmt.Sprintf("%s\n", obj.content))
	}

	for _, item := range obj.label {
		msg.WriteString(fmt.Sprintf("%s：%s\n", item[0], item[1]))
	}

	var data = map[string]string{
		"text": msg.String(),
	}

	marshal, _ := json.Marshal(data)
	return string(marshal)
}

func NewMessage() *Message {
	c := &Message{
		label: make([][]string, 0, 5),
	}
	return c
}
