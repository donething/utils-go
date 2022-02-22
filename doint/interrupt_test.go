package doint

import (
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	Init(func() {
		t.Logf("中断程序")
	})

	for {
		t.Logf("%d\n", time.Now().Unix())
		time.Sleep(1 * time.Second)
	}
}
