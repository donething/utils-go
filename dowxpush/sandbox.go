// Package dowxpush 微信测试号消息推送
package dowxpush

import (
	"encoding/json"
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"time"
)

const (
	// 获取 token
	tokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	// 推送消息
	pushURL = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s"
)

// Sandbox 微信测试号的数据
type Sandbox struct {
	appid   string           // 测试号的 appid
	secret  string           // 测试号的 secret
	token   string           // 根据 appid、secret 获取推送消息需要的 token
	expires time.Time        // token 的过期时间，以重复利用
	client  *dohttp.DoClient // 执行 http 请求
}

// NewSandbox 获取 Sandbox 实例，以推送消息
func NewSandbox(appid string, secret string) *Sandbox {
	client := dohttp.New(30*time.Second, false, false)
	return &Sandbox{appid: appid, secret: secret, token: "", expires: time.Now(), client: &client}
}

// 获取 token
func (s *Sandbox) getToken() error {
	// 如果 token 已存在且未过期，则继续使用，不重新获取
	if s.token != "" && time.Now().Before(s.expires) {
		return nil
	}

	// 需要重新获取
	bs, err := s.client.Get(fmt.Sprintf(tokenURL, s.appid, s.secret), nil)
	if err != nil {
		return err
	}
	// 解析并获取 token
	var token tokenResult
	err = json.Unmarshal(bs, &token)
	if err != nil {
		return err
	}
	if token.AccessToken == "" {
		return fmt.Errorf("未从文本中获取到 token：%s", string(bs))
	}

	// 获取 token 成功，设置 token 和 expires
	s.token = token.AccessToken
	duration, err := time.ParseDuration(fmt.Sprintf("%ds", token.ExpiresIn))
	if err != nil {
		return err
	}
	s.expires = time.Now().Add(duration)
	return nil
}

// PushTpl 推送模板消息
//
// payload 可以使用 GenGeneralTpl() 快速生成
//
// url 如果是有效链接，那么点击消息将会打开该链接
func (s *Sandbox) PushTpl(toUID string, tplID string, payload *map[string]interface{}, url string) error {
	// 获取、更新 token
	err := s.getToken()
	if err != nil {
		return fmt.Errorf("获取 token 出错：%w", err)
	}

	// 推送(post)的数据
	data := map[string]interface{}{"touser": toUID, "template_id": tplID,
		"url": url, "data": payload}
	bs, err := s.client.PostJSONObj(fmt.Sprintf(pushURL, s.token), data, nil)
	if err != nil {
		return fmt.Errorf("推送消息时网络出错：%w", err)
	}

	// 解析、判断推送的结果
	var result PushResult
	err = json.Unmarshal(bs, &result)
	if err != nil {
		return fmt.Errorf("解析推送响应 JSON 文本时出错：%w", err)
	}
	if result.Errcode != 0 {
		return fmt.Errorf("推送时出错：%s", string(bs))
	}

	return nil
}

// GenGeneralTpl 生成通用消息模板
func (s *Sandbox) GenGeneralTpl(title string, msg string, time string) *map[string]interface{} {
	return &map[string]interface{}{
		"title": map[string]string{"value": title + "\n"},
		"msg":   map[string]string{"value": msg + "\n\n"},
		"time":  map[string]string{"value": time},
	}
}
