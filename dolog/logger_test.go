package dolog

import (
	"testing"
	"time"
)

func TestInitLog(t *testing.T) {
	i, w, e, err := InitLog("./logger_test.log", 0, false)
	CkPanic(err)

	tick := time.Tick(3 * time.Second)
	for range tick {
		i.Println(time.Now().String())
		w.Println(time.Now().String())
		e.Println(time.Now().String())
	}
}
