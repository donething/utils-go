package dowxpush

// 获取 token 的结果
type tokenResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

// PushResult 推送消息的响应
type PushResult struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Msgid   int64  `json:"msgid"`
}
