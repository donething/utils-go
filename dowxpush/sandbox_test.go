package dowxpush

import (
	"log"
	"testing"
	"time"
)

var sandbox = NewSandbox("xxx", "xxx")

func TestSandbox_getToken(t *testing.T) {
	err := sandbox.getToken()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("token: %v\n", sandbox)
}

func TestSandbox_PushTpl(t *testing.T) {
	payload := sandbox.genGeneralTpl("测试标题", "测试消息内容", time.Now().String())
	err := sandbox.PushTpl("xxx", "xxx", payload, "")
	if err != nil {
		log.Fatal(err)
	}
}
