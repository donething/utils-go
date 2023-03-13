package dowx

import (
	"testing"
)

var sandbox = NewSandbox("xxx", "xxx")

func TestSandbox_PushTpl(t *testing.T) {
	err := sandbox.PushTpl("xxx", "xxx", "测试标题", "测试消息内容", "")
	if err != nil {
		t.Fatal(err)
	}
}
