package entity

type PushMsg struct {
	Touser     string      `json:"touser"`
	TemplateID string      `json:"template_id"`
	URL        string      `json:"url"`
	Data       interface{} `json:"data"`
}

type PushDataItem struct {
	Text  string `json:"text"`
	Color string `json:"color"`
}
