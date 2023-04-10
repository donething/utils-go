package dolive

import (
	"fmt"
	"testing"
)

var url = "https://cn-gddg-cu-01-04.bilivideo.com/live-bvc/777389/live_8739477_3713195_bluray.flv?expires=1681150168&pt=web&deadline=1681150168&len=0&oi=1885426465&platform=web&qn=10000&trid=1000aac0b115c9f448d58f6cfd28176b14f7&uipk=100&uipv=100&nbs=1&uparams=cdn,deadline,len,oi,platform,qn,trid,uipk,uipv,nbs&cdn=cn-gotcha01&upsig=0cb02b7251270d32e51915b15bcb3e34&sk=54c5c20ce223a922342f4414780ac3fe&p2p_type=0&sl=10&free_type=0&mid=10338719&sid=cn-gddg-cu-01-04&chash=0&sche=ban&score=18&pp=rtmp&source=one&trace=8a0&site=e35d773b7046f4632bc380ecfd566d1c&order=1"

func TestNew(t *testing.T) {
	live := New[string](url, "D:/Tmp/live/bili.flv", 10*1024*1024)

	err := live.Capture(BiliHeader, "数据", deal)
	if err != nil {
		panic(err)
	}
}

func deal(path string, data string) error {
	fmt.Printf("收到的数据：%s\n", data)
	fmt.Printf("处理视频：%s\n", path)
	return nil
}
