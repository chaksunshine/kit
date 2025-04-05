package authorized

import (
	"context"
	"errors"
	"fmt"
	"github.com/chaksunshine/kit/catLog"
	"github.com/chaksunshine/kit/thread"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Grpc请求信息
// @author fuzeyu
// @date 2025/4/5
type GrpcValidate struct {
	authorized
}

// 获取请求上下文
// @param timeout 超时时间
func (obj *GrpcValidate) CreateRequestContext(timeout ...int) context.Context {

	var seconds = 3
	if len(timeout) > 0 {
		seconds = timeout[0]
	}

	token, err := obj.newRequestToken()
	if err != nil {
		catLog.Error(fmt.Sprintf("生成Grpc签名密钥信息失败,直接创建空上下文 %s", err.Error()))
		return thread.CtxRequest(seconds)
	}
	md := metadata.Pairs(
		fieldsRequestToken, token,
	)
	return metadata.NewOutgoingContext(thread.CtxRequest(seconds), md)
}

// 验证请求参数
// @param ctx 上下文
// @param req 请求参数
// @param info 请求服务端信息
// @param handler 下一次回调
func (obj *GrpcValidate) ValidateRequest(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	incomingContext, b := metadata.FromIncomingContext(ctx)
	if b == false {
		return nil, ErrValidateAuthorizedFail
	}
	get := incomingContext.Get(fieldsRequestToken)
	if len(get) != 1 {
		return nil, ErrValidateAuthorizedFail
	}
	if err := obj.validateRequestToken(get[0]); err != nil {
		return nil, ErrValidateAuthorizedFail
	}

	// 调用后续的方法
	return handler(ctx, req)
}

// 创建消息流验证参数
// @param ctx 上下文
// @param ss 消息信息
// @param info 请求服务端信息
// @param handler 下一次回调
func (obj *GrpcValidate) ValidateRequestSteam(srv any, serverStream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	ctx := serverStream.Context()
	incomingContext, b := metadata.FromIncomingContext(ctx)
	if b == false {
		return ErrValidateAuthorizedFail
	}
	get := incomingContext.Get(fieldsRequestToken)
	if len(get) != 1 {
		return ErrValidateAuthorizedFail
	}
	if err := obj.validateRequestToken(get[0]); err != nil {
		return ErrValidateAuthorizedFail
	}

	// 调用后续的方法
	return handler(srv, serverStream)
}

// @param safeKey 安全密钥
func NewGrpcValidate(securityKey string) (*GrpcValidate, error) {
	c := &GrpcValidate{}
	if err := c.initSecurityKey(securityKey); err != nil {
		return nil, err
	}
	return c, nil
}

var ErrValidateAuthorizedFail = errors.New("验证请求中的签名信息失败")
