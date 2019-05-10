package entity

// 推送消息的模板
// 参数Data传递nil即可
type PushMsg struct {
	Touser     string              `json:"touser"`
	TemplateID string              `json:"template_id"`
	URL        string              `json:"url"`
	Data       map[string]dataItem `json:"data"`
}

// 模板中填充的数据
type dataItem struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

// 添加模板中的数据
// 参数key: 模板数据中变量的名字如"time.DATA"，则key为"time"
// 参数value: 模板数据值
// 参数color: 数据文本的颜色
func (p *PushMsg) AddData(key string, value string, color string) *PushMsg {
	if p.Data == nil {
		p.Data = make(map[string]dataItem)
	}
	p.Data[key] = dataItem{Value: value, Color: color}
	return p
}
