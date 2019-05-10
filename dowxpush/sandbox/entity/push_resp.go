package entity

// 微信沙盒消息推送后的返回信息
type PushResp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Msgid   int64  `json:"msgid"`
}
