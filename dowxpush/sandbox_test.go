package dowxpush

import (
	"testing"
	"time"
)

var sandbox = NewSandbox("xxx", "xxx")

func TestSandbox_PushTpl(t *testing.T) {
	payload := sandbox.GenGeneralTpl("测试标题", "测试消息内容", time.Now().String())
	err := sandbox.PushTpl("xxx", "xxx", payload, "")
	if err != nil {
		t.Fatal(err)
	}
}
