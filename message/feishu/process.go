package feishu

import (
	"fmt"
	"github.com/chaksunshine/kit/catLog"
	"github.com/chaksunshine/kit/thread"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"time"
)

// 发送消息
// @author fuzeyu
// @date 2025/4/2
type Process struct {
	notify chan *FeishuApp
}

// 发送消息
// @param app 配置信息
// @param message 消息信息
func (obj *Process) Send(message *Message, app ...*FeishuApp) {
	parser := message.parser()
	for _, item := range app {
		item.msg = parser

		obj.notify <- item
		time.Sleep(time.Millisecond * 20)
	}
}

// 发送消息
func (obj *Process) sendMessage(app *FeishuApp) {

	client := lark.NewClient(app.AppId, app.AppSecret)

	// 2. 构建消息请求
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			ReceiveId(app.UserOpenId).
			Content(app.msg).
			Build()).
		Build()

	// 3. 发送消息
	resp, err := client.Im.Message.Create(thread.CtxRequest(), req)
	if err != nil {
		catLog.Error(fmt.Sprintf("发送飞书消息失败 创建消息失败 %s", err.Error()))
		return
	}

	if !resp.Success() {
		catLog.Error(fmt.Sprintf("发送飞书消息失败 响应发送消息失败 [%d] %s %s", resp.Code,
			resp.Msg, resp.RequestId()))
		return
	}
}

// 开始消费
func (obj *Process) consume() {
	for app := range obj.notify {
		obj.sendMessage(app)
	}
}

// @param processNumber 发送消息的线程数
func NewProcess(processNumber int) *Process {
	c := &Process{
		notify: make(chan *FeishuApp, 20),
	}

	if processNumber <= 0 {
		catLog.FatalError(fmt.Sprintf("飞书工作线程数异常 %d", processNumber))
	}

	for index := 0; index < processNumber; index++ {
		go c.consume()
	}
	return c
}
