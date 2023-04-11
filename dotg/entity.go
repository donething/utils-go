package dotg

type IMedia interface {
	String()
}

// InputMedia 媒体类型
type InputMedia struct {
	// 媒体的类型。可选 TypeAudio、TypeDocument、TypePhoto、TypeVideo
	Type string `json:"type"`

	// 媒体内容。可以为 字节数组、URL
	Media interface{} `json:"media"`

	// 缩略图
	Thumbnail interface{} `json:"thumbnail,omitempty"`

	// 标题
	Caption string `json:"caption,omitempty"`

	// 标题的解析模式。推荐使用"MarkdownV2"
	ParseMode string `json:"parse_mode,omitempty"`

	// 是否遮罩（图片、视频）
	HasSpoiler bool `json:"has_spoiler,omitempty"`

	// 视频是否为流，为 true 可预览播放
	SupportsStreaming bool `json:"supports_streaming,omitempty"`
}

// Meida 的可选类型
const (
	// TypeAnimation Gif 或 无声视频
	TypeAnimation = "animation"
	TypeAudio     = "audio"
	TypeDocument  = "document"
	TypePhoto     = "photo"
	TypeVideo     = "video"
)

// 标题的解析模式
const (
	ParseMK2  = "MarkdownV2"
	ParseHTML = "HTML"
)
