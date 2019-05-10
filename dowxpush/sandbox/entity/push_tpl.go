package entity

// 推送消息的模板
type PushMsg struct {
	Touser     string      `json:"touser"`
	TemplateID string      `json:"template_id"`
	URL        string      `json:"url"`
	Data       interface{} `json:"data"`
}
type PushDataItem struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

// 推送的消息的数据
// 通用模板
type TplGeneral struct {
	Time  PushDataItem `json:"time"`
	Title PushDataItem `json:"title"`
	Msg   PushDataItem `json:"msg"`
}
