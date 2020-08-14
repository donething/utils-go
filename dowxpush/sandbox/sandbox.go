// 微信测试号，消息推送
package sandbox

import (
	"encoding/json"
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"github.com/donething/utils-go/dowxpush/sandbox/entity"
	"time"
)

const (
	pushWXURL    = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s"
	pushTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
)

type WXSendbox struct {
	pushURL string
	client  *dohttp.DoClient
}

// 创建微信测试号推送的实例
func NewSandbox(appID string, appSecret string) (wx WXSendbox, err error) {
	wx.client = dohttp.New(60*time.Second, false, false)
	tokenText, _, err := wx.client.GetText(fmt.Sprintf(pushTokenURL, appID, appSecret), nil)
	if err != nil {
		return
	}

	var tokenJson entity.PushToken
	err = json.Unmarshal([]byte(tokenText), &tokenJson)
	if err != nil {
		return wx, fmt.Errorf("unmarshal json err. error: %s, json text: %s", err, tokenText)
	}
	if tokenJson.AccessToken == "" {
		return wx, fmt.Errorf("access_token is blank: %s", tokenText)
	}
	wx.pushURL = fmt.Sprintf(pushWXURL, tokenJson.AccessToken)
	return
}

// 推送模板消息
func (wx *WXSendbox) PushTpl(pushMsg entity.PushMsg) (resp entity.PushResp, err error) {
	// 创建POST的json数据
	bs, err := json.Marshal(pushMsg)
	if err != nil {
		return
	}
	jsonText := string(bs)
	// 推送消息
	bs, _, err = wx.client.PostJSONString(wx.pushURL, jsonText, nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(bs, &resp)
	return
}
