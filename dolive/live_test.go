package dolive

import (
	"testing"
)

var url = "https://cn-gddg-cu-01-05.bilivideo.com/live-bvc/458532/live_474595627_50840492_bluray.flv?expires=1681006135&pt=web&deadline=1681006135&len=0&oi=1885426465&platform=web&qn=10000&trid=1000d13a1ccf94d749009f95d50766d2adce&uipk=100&uipv=100&nbs=1&uparams=cdn,deadline,len,oi,platform,qn,trid,uipk,uipv,nbs&cdn=cn-gotcha01&upsig=4fb8f741552a1d314d3352f0fabe9ec2&sk=b6df54012328e502455d040433e3ec02&p2p_type=0&sl=3&free_type=0&mid=10338719&sid=cn-gddg-cu-01-05&chash=1&sche=ban&score=16&pp=rtmp&source=one&trace=8a0&site=555cbf203b4a917f9c57ef9e40b5640e&order=1"

func TestNew(t *testing.T) {
	live := New(url, "D:/Tmp/live/bili.flv", 0)

	err := live.Capture(BiliHeader)
	if err != nil {
		panic(err)
	}
}
