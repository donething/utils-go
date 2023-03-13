package dotg

const (
	// TypeAudio 媒体的类型
	TypeAudio    = "audio"
	TypeDocument = "document"
	TypePhoto    = "photo"
	TypeVideo    = "video"
)

// Media 媒体文件
type Media struct {
	// 媒体的类型。可选 TypeAudio、TypeDocument、TypePhoto、TypeVideo
	Type string `json:"type"`
	// 媒体内容。可以为 字节数组、URL
	Media interface{} `json:"media"`
	// 标题
	Caption string `json:"caption"`
	// 标题的解析模式。推荐使用"MarkdownV2"
	ParseMode string `json:"parse_mode"`
}

// Message 为发送消息后返回的响应
//
// OK 为 true 表示成功，false 为失败
type Message struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description，omitempty"`
}
