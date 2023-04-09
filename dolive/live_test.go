package dolive

import (
	"fmt"
	"testing"
)

var url = "https://cn-gddg-cu-01-05.bilivideo.com/live-bvc/789977/live_474595627_50840492_bluray.flv?expires=1681009962&pt=web&deadline=1681009962&len=0&oi=1885426465&platform=web&qn=10000&trid=100081aa6a1f409d4294b9ca10029cd70cf6&uipk=100&uipv=100&nbs=1&uparams=cdn,deadline,len,oi,platform,qn,trid,uipk,uipv,nbs&cdn=cn-gotcha01&upsig=0f30e893f49402e673594e07d5897646&sk=b6df54012328e502455d040433e3ec02&p2p_type=0&sl=3&free_type=0&mid=10338719&sid=cn-gddg-cu-01-05&chash=1&sche=ban&score=16&pp=rtmp&source=one&trace=8a0&site=01addb68232f0230098b083cef0b31de&order=1"

func TestNew(t *testing.T) {
	live := New(url, "D:/Tmp/live/bili.flv", 10*1024*1024)

	err := live.Capture(BiliHeader, deal)
	if err != nil {
		panic(err)
	}
}

func deal(path string) error {
	fmt.Printf("处理视频：%s\n", path)
	return nil
}
