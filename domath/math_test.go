package domath

import (
	"testing"
	"time"
)

func TestRandInt(t *testing.T) {
	for i := 0; i <= 100; i++ {
		num := RandInt(10, 20)
		// 需要等待一会，避免连续获取同一个的随机数
		time.Sleep(10 * time.Millisecond)
		t.Log(num)
	}
}
