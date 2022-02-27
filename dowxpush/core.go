package dowxpush

import (
	"encoding/json"
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"time"
)

// Core 微信消息推送的核心信息
type Core struct {
	appid   string    // 测试号的 appid
	secret  string    // 测试号的 secret
	token   string    // 根据 appid、secret 获取推送消息需要的 token
	expires time.Time // token 的过期时间，以重复利用
}

var client = dohttp.New(30*time.Second, false, false)

// 获取 token
func (c *Core) getToken(url string) error {
	// 如果 token 已存在且未过期，则继续使用，不重新获取
	if c.token != "" && time.Now().Before(c.expires) {
		return nil
	}

	// 需要重新获取
	bs, err := client.Get(url, nil)
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
		return fmt.Errorf("无法从文本中获取到 token: %s", string(bs))
	}

	// 获取 token 成功，设置 token 和 expires
	c.token = token.AccessToken
	// 将 token 过期时间减少 3 分钟，以容错
	duration, err := time.ParseDuration(fmt.Sprintf("%ds", token.ExpiresIn-180))
	if err != nil {
		return err
	}
	c.expires = time.Now().Add(duration)
	return nil
}

// 推送消息
func (c *Core) push(tokenURL string, sendURL string, data interface{}) error {
	// 获取、更新 token
	err := c.getToken(fmt.Sprintf(tokenURL, c.appid, c.secret))
	if err != nil {
		return fmt.Errorf("获取 token 出错：%w", err)
	}

	// 推送(post)的数据
	bs, err := client.PostJSONObj(fmt.Sprintf(sendURL, c.token), data, nil)
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
