package dowxpush

import (
	"testing"
)

var aid = 111 // 应用 ID
var qy = NewQiYe("xxx", "xxx")

func TestQiYe_PushText(t *testing.T) {
	err := qy.PushText(aid, "测试文本消息，不错", "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestQiYe_PushCard(t *testing.T) {
	err := qy.PushCard(aid, "消息标题", "测试文本消息，不错"+
		GenHyperlink("https://developer.work.weixin.qq.com/document/path/90236", "微信开发者中心"),
		"", "https://www.jianshu.com/p/182ea14af3f2", "打开")
	if err != nil {
		t.Fatal(err)
	}
}

func TestQiYe_PushMarkdown(t *testing.T) {
	err := qy.PushMarkdown(aid, `您的会议室已经预定，稍后会同步到**邮箱**`+
		GenMdInfoText("测试信息文本DIV"), "")
	if err != nil {
		t.Fatal(err)
	}
}
