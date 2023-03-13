package dowx

// 请求的相应

// 获取 token 的结果
type tokenResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

// PushResult 推送消息的响应
type PushResult struct {
	// 通用部分
	Errcode int         `json:"errcode"`
	Errmsg  string      `json:"errmsg"`
	Msgid   interface{} `json:"msgid"` // 消息 ID，企业消息返回 数字，测试号消息返回 字符串

	// 企业微信部分
	Invaliduser  string `json:"invaliduser"`
	Invalidparty string `json:"invalidparty"`
	Invalidtag   string `json:"invalidtag"`
	ResponseCode string `json:"response_code"`
}

// 微信测试号消息推送

// SBMsgItem 微信测试号消息的项
type SBMsgItem struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

// SBMsg 微信测试号消息
type SBMsg struct {
	Title struct {
		Value string `json:"value"`
		Color string `json:"color"`
	} `json:"title"`
	Msg struct {
		Value string `json:"value"`
		Color string `json:"color"`
	} `json:"msg"`
	Time struct {
		Value string `json:"value"`
		Color string `json:"color"`
	} `json:"time"`
}

// 企业微信消息推送

// QYMsg 企业消息模板
//
// 不完整，实际使用前要自己添加消息内容。可用 QYMsgText、QYMsgCard
type QYMsg struct {
	Touser               string `json:"touser"` // 推送目标用户，当为"@all"时推送所有人
	Toparty              string `json:"toparty"`
	Totag                string `json:"totag"`
	Msgtype              string `json:"msgtype"` // 消息类型，如"text"、"textcard"
	Agentid              int    `json:"agentid"` // 应用 ID
	Safe                 int    `json:"safe"`
	EnableIdTrans        int    `json:"enable_id_trans"`
	EnableDuplicateCheck int    `json:"enable_duplicate_check"`
}

// QYMsgItemText 企业文本消息
type QYMsgItemText struct {
	Content string `json:"content"`
}

// QYMsgText 企业文本消息
type QYMsgText struct {
	*QYMsg
	Text QYMsgItemText `json:"text"`
}

// QYMsgItemCard 企业卡片消息
type QYMsgItemCard struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`    // 点击卡片内容跳转链接
	Btntxt      string `json:"btntxt"` // 可跳转链接的文本
}

// QYMsgCard 企业卡片消息
type QYMsgCard struct {
	*QYMsg
	Textcard QYMsgItemCard `json:"textcard"`
}

// QYMsgMarkdown 企业 Markdown 消息
type QYMsgMarkdown struct {
	*QYMsg
	Markdown QYMsgItemText `json:"markdown"`
}
