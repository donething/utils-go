package dotg

// InputMedia 媒体的数据
type InputMedia struct {
	// 媒体的类型。可选 TypeAudio、TypeDocument、TypePhoto、TypeVideo
	Type string `json:"type"`

	// 媒体内容。为 *io.Reader（可读流）、string（本地文件地址，如"file:///homevideo.mp4"）
	Media interface{} `json:"media"`

	// 缩略图。始终传递 io.Reader 类型，但发送时读取到表单后，需要设置为字符串("attach://thumb1.jpg")指向表单
	Thumbnail interface{} `json:"thumbnail,omitempty"`

	// 标题
	Caption string `json:"caption,omitempty"`

	// 标题的解析模式。推荐使用"MarkdownV2"
	ParseMode string `json:"parse_mode,omitempty"`

	// 视频是否为流，为 true 可预览播放
	SupportsStreaming bool `json:"supports_streaming,omitempty"`

	// 分辨率。如果设置，TG 在播放、在聊天框中显示时，宽度、高度可能不准确
	Width  int `json:"width"`
	Height int `json:"height"`

	// 是否遮罩（图片、视频）
	HasSpoiler bool `json:"has_spoiler,omitempty"`

	// 	非 TG 属性
	//  文件名
	Name string
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
