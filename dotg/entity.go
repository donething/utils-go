package dotg

import "io"

// MediaData 媒体的数据
// 引入单独的 MediaData 是因为 SendMediaGroup() 时要将 InputMedia 中的 Media 转为字符串（指向 multipart form 中的该文件）
// 即 Media 要为 *io.Reader 和 string 两种类型，未免传错参数，不定以为 interface{}，而分拆属性
type MediaData struct {
	// 媒体的类型。可选 TypeAudio、TypeDocument、TypePhoto、TypeVideo
	Type string `json:"type"`

	// 标题
	Caption string `json:"caption,omitempty"`

	// 标题的解析模式。推荐使用"MarkdownV2"
	ParseMode string `json:"parse_mode,omitempty"`

	// 是否遮罩（图片、视频）
	HasSpoiler bool `json:"has_spoiler,omitempty"`

	// 视频是否为流，为 true 可预览播放
	SupportsStreaming bool `json:"supports_streaming,omitempty"`

	// 分辨率。如果设置，TG 在播放、在聊天框中显示时，宽度、高度可能不准确
	Width  int `json:"width"`
	Height int `json:"height"`

	// 	非 TG 属性
	//  文件名
	Name string
}

// MediaForm 发送的媒体的表单数据
type MediaForm struct {
	*MediaData

	// 指向的媒体，为 multipart form 中该文件的指向。
	// 如 "attach://1"，表示该 MediaData 为文件表单中的第一个文件的信息
	Media string `json:"media"`

	// 指向的缩略图。意义和 Media 相同
	Thumbnail string `json:"thumbnail,omitempty"`
}

// InputMedia 媒体类型
type InputMedia struct {
	*MediaData

	// 媒体内容。为 字节数组、字符串、文件的 *Reader
	Media io.Reader `json:"media"`

	// 缩略图
	Thumbnail io.Reader `json:"thumbnail,omitempty"`
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
