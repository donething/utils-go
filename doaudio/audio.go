package doaudio

import (
	"bytes"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"io"
	"io/ioutil"
)

const (
	BufSize8192 = 8192
)

// 如果已有音频正在播放，则关闭后开始播放新音频
var ctx *oto.Context

// 播放音频
// 参数audio: 音频文件的字节码
// 参数channelNum: 单声道（1）或立体声（2）
// 参数bitDepthInBytes: 只能选择1或2，大多数情况为2
func PlayAudio(audio []byte, channelNum int, bitDepthInBytes int, bufferSizeInBytes int) (err error) {
	if ctx != nil {
		err = ctx.Close()
		if err != nil {
			return
		}
	}

	d, err := mp3.NewDecoder(ioutil.NopCloser(bytes.NewReader(audio)))
	if err != nil {
		return
	}

	ctx, err = oto.NewContext(d.SampleRate(), channelNum, bitDepthInBytes, bufferSizeInBytes)
	if err != nil {
		return
	}
	// 需要结束后close Context，否则会报错：oto: NewContext can be called only once
	defer ctx.Close()

	p := ctx.NewPlayer()
	defer p.Close()

	if _, err = io.Copy(p, d); err != nil {
		return
	}

	ctx = nil
	return
}
