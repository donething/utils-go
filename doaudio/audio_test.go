package doaudio

import (
	"github.com/donething/utils-go/doaudio/audios"
	"log"
	"sync"
	"testing"
)

func Test_PlayAudio(t *testing.T) {
	var wg sync.WaitGroup

	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := PlayAudio(audios.Beep, 2, 2, BufSize8192)
			if err != nil {
				log.Printf("error: %v\n", err)
			}
			t.Logf("第%d次测试", i)
		}(i)
	}
	wg.Wait()
}

func Benchmark_PlayAudio(b *testing.B) {
	for i := 0; i < b.N; i++ { //use b.N for looping
		err := PlayAudio(audios.Success, 2, 2, BufSize8192)
		if err != nil {
			log.Printf("error: %v\n", err)
		}
	}
}
