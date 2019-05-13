package doaudio

import (
	"github.com/donething/utils-go/doaudio/audios"
	"log"
	"testing"
)

func Test_PlayAudio(t *testing.T) {
	for i := 1; i <= 5; i++ {
		err := PlayAudio(audios.Beep, 2, 2, BufSize8192)
		if err != nil {
			log.Printf("error: %v\n", err)
		}
	}
}

func Benchmark_PlayAudio(b *testing.B) {
	for i := 0; i < b.N; i++ { //use b.N for looping
		err := PlayAudio(audios.Success, 2, 2, BufSize8192)
		if err != nil {
			log.Printf("error: %v\n", err)
		}
	}
}
