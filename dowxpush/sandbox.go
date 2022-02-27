// Package dowxpush 微信测试号消息推送
package dowxpush

import (
	"github.com/donething/utils-go/dotext"
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

// Push 推送消息
//
// url 如果是有效链接，那么点击消息将会打开该链接
func (s *Sandbox) Push(toUID string, tplID string, payload interface{}, url string) error {
	// 推送(POST)的数据
	data := map[string]interface{}{"touser": toUID, "template_id": tplID,
		"url": url, "data": payload}

	return s.Core.push(sbTokenURL, sbSendURL, data)
}

// PushTpl 推送模板消息
//
// url 如果是有效链接，那么点击消息将会打开该链接
func (s *Sandbox) PushTpl(toUID string, tplID string, title string, msg string, url string) error {
	payload := &SBMsg{
		Title: SBMsgItem{Value: title + "\n"},
		Msg:   SBMsgItem{Value: msg + "\n"},
		Time:  SBMsgItem{Value: dotext.FormatDate(time.Now(), dotext.TimeFormat)},
	}
	// 推送(POST)的数据
	data := map[string]interface{}{"touser": toUID, "template_id": tplID,
		"url": url, "data": payload}

	return s.Core.push(sbTokenURL, sbSendURL, data)
}
