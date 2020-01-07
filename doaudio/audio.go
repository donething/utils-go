package doaudio

import (
	"bytes"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"io"
	"io/ioutil"
	"sync"
)

const (
	BufSize8192 = 8192
)

var ctx *oto.Context

// 仅调用一次oto.NewContext()，已避免"panic: oto: NewContext can be called only once"
// 参考：https://github.com/hajimehoshi/go-mp3/issues/28#issuecomment-523265030
var once sync.Once

// 播放音频
// 参数audio: 音频文件的字节码
// 参数channelNum: 单声道（1）或立体声（2）
// 参数bitDepthInBytes: 只能选择1或2，大多数情况为2
func PlayAudio(audio []byte, channelNum int, bitDepthInBytes int, bufferSizeInBytes int) (err error) {
	d, err := mp3.NewDecoder(ioutil.NopCloser(bytes.NewReader(audio)))
	if err != nil {
		return
	}

	once.Do(func() {
		ctx, err = oto.NewContext(d.SampleRate(), channelNum, bitDepthInBytes, bufferSizeInBytes)
		if err != nil {
			panic(err)
		}
	})

	p := ctx.NewPlayer()
	defer p.Close()

	if _, err = io.Copy(p, d); err != nil {
		return
	}

	return
}
