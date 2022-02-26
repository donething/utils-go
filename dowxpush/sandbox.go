// Package dowxpush 微信测试号消息推送
package dowxpush

import (
	"time"
)

const (
	// 获取 Token
	sbTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	// 发送消息
	sbSendURL = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s"
)

// Sandbox 微信测试号消息推送对象
type Sandbox struct {
	*Core
}

// NewSandbox 获取 Sandbox 实例，以推送消息
func NewSandbox(appid string, secret string) *Sandbox {
	return &Sandbox{
		Core: &Core{appid: appid, secret: secret, token: "", expires: time.Now()},
	}
}

// PushTpl 推送模板消息
//
// payload 可以使用 GenGeneralTpl() 快速生成
//
// url 如果是有效链接，那么点击消息将会打开该链接
func (s *Sandbox) PushTpl(toUID string, tplID string, payload interface{}, url string) error {
	// 推送(POST)的数据
	data := map[string]interface{}{"touser": toUID, "template_id": tplID,
		"url": url, "data": payload}

	return s.Core.Push(sbTokenURL, sbSendURL, data)
}

// GenGeneralTpl 生成通用消息模板
func (s *Sandbox) GenGeneralTpl(title string, msg string, time string) *SBMsg {
	return &SBMsg{
		Title: struct {
			Value string `json:"value"`
			Color string `json:"color"`
		}{Value: title},
		Msg: struct {
			Value string `json:"value"`
			Color string `json:"color"`
		}{Value: msg},
		Time: struct {
			Value string `json:"value"`
			Color string `json:"color"`
		}{Value: time},
	}
}
