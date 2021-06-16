// Package dotgpush 绑定机器人的 token 后即可发送消息
// tg := dotgpush.NewTGBot(token)
// tg.SendMessage(chatID, text)
package dotgpush

import (
	"encoding/json"
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"net/url"
	"time"
)

// TGBot TG 推送发送消息的机器人实例
type TGBot struct {
	token string
}

const (
	// API
	urlSendMsg        = "https://api.telegram.org/%s/sendMessage"
	urlSendMediaGroup = "https://api.telegram.org/%s/sendMediaGroup"

	// Audio 媒体的类型
	Audio    = "audio"
	Document = "document"
	Photo    = "photo"
	Video    = "video"
)

var (
	client = dohttp.New(30*time.Second, false, false)
)

// NewTGBot 创建新的 Telegram 推送机器人
//
// token 新机器人的 token
func NewTGBot(token string) *TGBot {
	return &TGBot{token: "bot" + token}
}

// SetProxy 设置网络代理
// 格式参考 dohttp.ProxySocks5、dohttp.ProxyHttp
func (bot *TGBot) SetProxy(proxyStr string) {
	client.SetProxy(proxyStr)
}

// SendMessage 将文本消息发送到频道
func (bot *TGBot) SendMessage(chatID string, text string) (*Message, error) {
	form := url.Values{
		"chat_id":    []string{chatID},
		"text":       []string{text},
		"parse_mode": []string{"markdown"},
	}
	bs, err := client.PostForm(fmt.Sprintf(urlSendMsg, bot.token), form.Encode(), nil)
	if err != nil {
		return nil, err
	}
	// 返回的消息
	var msg Message
	err = json.Unmarshal(bs, &msg)
	return &msg, err
}

// SendMediaGroup 发送一个媒体集
func (bot *TGBot) SendMediaGroup(chatID string, album []Media) (*Message, error) {
	form := make(map[string]string)
	form["chat_id"] = chatID

	// 存放媒体文件的数组，将被发送
	filesList := make(map[string]interface{})
	for i, m := range album {
		// 当文件形式为二进制数组数据时，需要在表单中加入该数据的指向标志
		if bs, ok := m.Media.([]byte); ok {
			album[i].Media = fmt.Sprintf("attach://%d", i)
			filesList[fmt.Sprintf("%d", i)] = bs
		}
	}

	// 将 album 序列化后作为表单"media"的值发送
	mediaStr, _ := json.Marshal(album)
	form["media"] = string(mediaStr)

	bs, err := client.PostFiles(fmt.Sprintf(urlSendMediaGroup, bot.token), filesList, form, nil)
	if err != nil {
		return nil, err
	}

	// 返回的消息
	var msg Message
	err = json.Unmarshal(bs, &msg)
	return &msg, err
}
