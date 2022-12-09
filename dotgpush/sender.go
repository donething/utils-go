// Package dotgpush 绑定机器人的 token 后即可发送消息
// tg := dotgpush.NewTGBot(token)
// tg.SendMessage(chatID, text)
package dotgpush

import (
	"encoding/json"
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"net/url"
	"strconv"
	"strings"
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
	client = dohttp.New(false, false)
)

// NewTGBot 创建新的 Telegram 推送机器人
//
// token 新机器人的 token
func NewTGBot(token string) *TGBot {
	return &TGBot{token: "bot" + token}
}

// SetProxy 设置网络代理
// 格式参考 dohttp.ProxySocks5、dohttp.ProxyHttp
func (bot *TGBot) SetProxy(proxyStr string) error {
	return client.SetProxy(proxyStr)
}

// SendMessage 发送Markdown文本消息
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
//
// 当只设置第一个媒体元素的`Caption`时，将作为该集的标题来显示
func (bot *TGBot) SendMediaGroup(chatID string, album []Media) (*Message, error) {
	// 避免改动原数据，以免重试发送时丢失二进制文件数据
	var newAlbum = make([]Media, 0, len(album))

	// 其它属性
	form := make(map[string]string)
	form["chat_id"] = chatID

	// 当album为本地图片文件（二进制数据），需要作为文件发送
	filesList := make(map[string]interface{})
	for i, m := range album {
		n := Media{
			Type:    m.Type,
			Media:   m.Media,
			Caption: m.Caption,
		}

		// 此时需要在表单中加入该数据的指向标志
		if bs, ok := n.Media.([]byte); ok {
			n.Media = fmt.Sprintf("attach://%d", i)
			filesList[fmt.Sprintf("%d", i)] = bs
		}

		newAlbum = append(newAlbum, n)
	}

	// 将 album 序列化后作为表单"media"的值发送
	mediaStr, _ := json.Marshal(newAlbum)
	form["media"] = string(mediaStr)

	bs, err := client.PostFiles(fmt.Sprintf(urlSendMediaGroup, bot.token), filesList, form, nil)
	if err != nil {
		return nil, err
	}

	// 返回的消息
	var msg Message
	err = json.Unmarshal(bs, &msg)
	if err != nil {
		return nil, err
	}

	// 速率限制
	if msg.ErrorCode == 429 {
		secStr := msg.Description[strings.LastIndex(msg.Description, " ")+1:]
		sec, err := strconv.Atoi(secStr)
		if err != nil {
			return nil, fmt.Errorf("解析速率限制的等待时长时出错：%w", err)
		}

		fmt.Printf("由于速率限制，等待 %d 秒后重新发送\n", sec)
		time.Sleep(time.Duration(sec+1) * time.Second)
		return bot.SendMediaGroup(chatID, album)
	}

	return &msg, nil
}
