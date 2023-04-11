package dotg

// Message 为发送消息后返回的响应
//
// OK 为 true 表示成功，false 为失败
type Message struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description，omitempty"`
}
