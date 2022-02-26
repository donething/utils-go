package dowxpush

import (
	"encoding/json"
	"fmt"
	"github.com/donething/utils-go/dotext"
	"time"
)

const (
	// 获取 Token
	qyTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	// 发送消息
	qySendURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
)

// QiYe 企业微信消息推送对象
type QiYe struct {
	*Core
}

// NewQiYe 获取 QiYe 实例，以推送消息
func NewQiYe(corpid string, corpsecret string) *QiYe {
	return &QiYe{&Core{appid: corpid, secret: corpsecret, token: "", expires: time.Now()}}
}

// Push 推送消息
func (q *QiYe) Push(data interface{}) error {
	return q.Core.Push(qyTokenURL, qySendURL, data)
}

// PushText 推送文本消息
//
// agentid 应用 ID；users 推送的目标（多个以"|"分隔），为空表示推送到所有人
func (q *QiYe) PushText(agentid int, content string, users string) error {
	if users == "" {
		users = "@all"
	}

	data := QYMsgText{
		QYMsg: &QYMsg{
			Touser:  users,
			Msgtype: "text",
			Agentid: agentid,
		},
		Text: QYMsgItemText{
			Content: content + "\n" + dotext.FormatDate(time.Now(), dotext.TimeFormat),
		},
	}

	return q.Core.Push(qyTokenURL, qySendURL, data)
}

// PushCard 推送卡片消息
//
// agentid 应用 ID；users 推送的目标（多个以"|"分隔），为空表示推送到所有人
//
// description 可以用"\n"换行，可用 DIV 标签设置字体颜色。已内置 3 种 class 可直接使用：
// 灰色(gray)、高亮(highlight)、默认黑色(normal)
func (q *QiYe) PushCard(agentid int, title string, description string, users string,
	url string, btnTxt string) error {
	if users == "" {
		users = "@all"
	}

	data := QYMsgCard{
		QYMsg: &QYMsg{
			Touser:  users,
			Msgtype: "textcard",
			Agentid: agentid,
		},
		Textcard: QYMsgItemCard{
			Title:       title,
			Description: aTime() + "\n" + description,
			Url:         url,
			Btntxt:      btnTxt,
		},
	}

	return q.Core.Push(qyTokenURL, qySendURL, data)
}

// PushMarkdown 推送 Markdown 消息（目前非企业微信不支持该类型）
//
// agentid 应用 ID；users 推送的目标（多个以"|"分隔），为空表示推送到所有人
//
// @see [企业微信报警中关于markdown的用法 - 三度](https://www.cnblogs.com/sanduzxcvbnm/p/14266180.html)
func (q *QiYe) PushMarkdown(agentid int, content string, users string) error {
	if users == "" {
		users = "@all"
	}

	data := QYMsgMarkdown{
		QYMsg: &QYMsg{
			Touser:  users,
			Msgtype: "markdown",
			Agentid: agentid,
		},
		Markdown: QYMsgItemText{
			Content: content + "\n" + dotext.FormatDate(time.Now(), dotext.TimeFormat),
		},
	}

	bsData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("反序列出错：%s\n", err)
		return err
	}
	fmt.Printf("POST 数据：%s\n", string(bsData))
	return q.Core.Push(qyTokenURL, qySendURL, data)
}

// MdInfoText 生成 Markdown 中显示为绿色的文本
func (q *QiYe) MdInfoText(text string) string {
	return fmt.Sprintf("<font color='info'>%s</font>", text)
}

// MdCommentText 生成 Markdown 中显示为灰色的文本
func (q *QiYe) MdCommentText(text string) string {
	return fmt.Sprintf("<font color='comment'>%s</font>", text)
}

// MdWarningText 生成 Markdown 中显示为橙红色的文本
func (q *QiYe) MdWarningText(text string) string {
	return fmt.Sprintf("<font color='warning'>%s</font>", text)
}
