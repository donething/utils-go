package dowx

import (
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
	return q.Core.push(qyTokenURL, qySendURL, data)
}

// PushText 推送文本消息
//
// agentid 应用 ID
//
// content 消息内容，支持换行"\n"、以及超链接"A"
//
// users 推送的目标（多个以"|"分隔），为空表示推送到所有人
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
			Content: content,
		},
	}

	return q.Core.push(qyTokenURL, qySendURL, data)
}

// PushTextMsg 推送文本消息（含标题、正文、当前时间）
//
// agentid 应用 ID
//
// title 标题
//
// msg 消息内容，支持换行"\n"、以及超链接"A"
//
// users 推送的目标（多个以"|"分隔），为空表示推送到所有人
func (q *QiYe) PushTextMsg(agentid int, title string, msg string, users string) error {
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
			Content: title + "\n\n" + msg + "\n\n" + dotext.FormatDate(time.Now(), dotext.TimeFormat),
		},
	}

	return q.Core.push(qyTokenURL, qySendURL, data)
}

// PushCard 推送卡片消息
//
// agentid 应用 ID
//
// title 标题
//
// description 内容。可以用"\n"换行，可用设置部分字体颜色（已提供函数快速生成），不可含超链接
//
// users 推送的目标（多个以"|"分隔），为空表示推送到所有人
//
// url 跳转链接 由于不能为空""，当传递""时，将设为默认值 "https://example.com"
//
// btnTxt 跳转标识文本（仅在企业微信中有效，在微信中无效）
func (q *QiYe) PushCard(agentid int, title string, description string, users string,
	url string, btnTxt string) error {
	if users == "" {
		users = "@all"
	}
	if url == "" {
		url = "https://example.com"
	}

	data := QYMsgCard{
		QYMsg: &QYMsg{
			Touser:  users,
			Msgtype: "textcard",
			Agentid: agentid,
		},
		Textcard: QYMsgItemCard{
			Title: title,
			Description: GenCardGrayText(dotext.FormatDate(time.Now(), dotext.TimeFormat)) +
				"\n" + description,
			Url:    url,
			Btntxt: btnTxt,
		},
	}

	return q.Core.push(qyTokenURL, qySendURL, data)
}

// PushMarkdown 推送 Markdown 消息（目前非企业微信不支持该类型）
//
// agentid 应用 ID
//
// content 目前仅支持 Markdown 语法的子集
//
// users 推送的目标（多个以"|"分隔），为空表示推送到所有人
//
// @see https://developer.work.weixin.qq.com/document/path/90236#markdown%E6%B6%88%E6%81%AF
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
			Content: content + "\n\n" + dotext.FormatDate(time.Now(), dotext.TimeFormat),
		},
	}

	return q.Core.push(qyTokenURL, qySendURL, data)
}

// 快速生成器

// GenCardGrayText 生成注释（灰色）文本，仅适用于卡片消息
func GenCardGrayText(text string) string {
	return fmt.Sprintf("<div class='gray'>%s</div>", text)
}

// GenCardNormalText 生成普通（黑色）文本，仅适用于卡片消息
func GenCardNormalText(text string) string {
	return fmt.Sprintf("<div class='normal'>%s</div>", text)
}

// GenCardHighlightText 生成高亮（橙红色）文本，仅适用于卡片消息
func GenCardHighlightText(text string) string {
	return fmt.Sprintf("<div class='highlight'>%s</div>", text)
}

// GenMdInfoText 生成信息（绿色）文本，仅适用于 Markdown
func GenMdInfoText(text string) string {
	return fmt.Sprintf("<font color='info'>%s</font>", text)
}

// GenMdCommentText 生成注释（灰色）文本，仅适用于 Markdown
func GenMdCommentText(text string) string {
	return fmt.Sprintf("<font color='comment'>%s</font>", text)
}

// GenMdWarningText 生成警告（橙红色）文本，仅适用于 Markdown
func GenMdWarningText(text string) string {
	return fmt.Sprintf("<font color='warning'>%s</font>", text)
}

// GenHyperlink 生成超链接
func GenHyperlink(url string, title string) string {
	return fmt.Sprintf("<a href='%s'>%s</a>", url, title)
}
