package domath

import (
	"math/rand"
	"time"
)

// RandInt 获取 [min, max] 之间的随机整数
//
// 不要瞬间获取多个随机数（用 for 连接获取），此时会得到很多相同的随机数，可以间隔几毫秒
func RandInt(min, max int) int {
	// 创建一个随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 获取[min,max]之间的随机整数
	return r.Intn(max-min+1) + min
}
