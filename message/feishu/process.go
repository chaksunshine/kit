package feishu

import (
	"fmt"
	"github.com/chaksunshine/kit/catLog"
	"github.com/chaksunshine/kit/thread"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// 发送消息
// @author fuzeyu
// @date 2025/4/2
type process struct {
	notify chan *FeishuApp
}

// 发送消息
// @param app 配置信息
// @param message 消息信息
func (obj *process) Send(app *FeishuApp, message *Message) {
	app.msg = message.parser()
	obj.notify <- app
}

// 发送消息
func (obj *process) sendMessage(app *FeishuApp) {

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
func (obj *process) consume() {
	for app := range obj.notify {
		obj.sendMessage(app)
	}
}

func newProcess() *process {
	c := &process{
		notify: make(chan *FeishuApp, 20),
	}

	for index := 0; index < 3; index++ {
		go c.consume()
	}
	return c
}

var ProcessLogic = newProcess()
