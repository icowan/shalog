/**
 * @Time: 2020/5/1 15:43
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package wechat

import (
	"context"
	"github.com/chanxuehong/wechat/mp/message/callback/response"
	"github.com/go-kit/kit/log/level"

	"github.com/chanxuehong/wechat/mp/core"
	"github.com/chanxuehong/wechat/mp/menu"
	"github.com/chanxuehong/wechat/mp/message/callback/request"
	"github.com/go-kit/kit/log"
)

type Service interface {
	Callback(ctx context.Context) *core.Server
}

type service struct {
	logger log.Logger
	server *core.Server
}

func (s *service) Callback(ctx context.Context) *core.Server {
	return s.server
}

func (s *service) defaultMsgHandler(ctx *core.Context) {
	_ = level.Debug(s.logger).Log("default", "msg")
	_ = ctx.NoneResponse()
}

func (s *service) defaultEventHandler(ctx *core.Context) {
	_ = ctx.NoneResponse()
}

func (s *service) textMsgHandler(ctx *core.Context) {
	msg := request.GetText(ctx.MixedMsg)

	_ = level.Debug(s.logger).Log("fromUserName", msg.FromUserName, "toUserName", msg.ToUserName, "createTime", msg.CreateTime, "content", msg.Content)

	resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.CreateTime, msg.Content)
	//ctx.RawResponse(resp) // 明文回复
	_ = ctx.AESResponse(resp, 0, "什么？", nil) // aes密文回复
}

func (s *service) menuClickEventHandler(ctx *core.Context) {
	event := menu.GetClickEvent(ctx.MixedMsg)
	resp := response.NewText(event.FromUserName, event.ToUserName, event.CreateTime, "收到 click 类型的事件")
	//ctx.RawResponse(resp) // 明文回复
	_ = ctx.AESResponse(resp, 0, "点了一下", nil) // aes密文回复
}

func NewService(logger log.Logger, oriId, appId, token, base64AESKey string) Service {
	s := &service{logger: logger}
	mux := core.NewServeMux()
	mux.DefaultMsgHandleFunc(s.defaultMsgHandler)
	mux.DefaultEventHandleFunc(s.defaultEventHandler)
	mux.MsgHandleFunc(request.MsgTypeText, s.textMsgHandler)
	mux.EventHandleFunc(menu.EventTypeClick, s.menuClickEventHandler)
	s.server = core.NewServer(oriId, appId, token, base64AESKey, mux, nil)
	return s
}
