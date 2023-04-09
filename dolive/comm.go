package dolive

// BiliHeader 哔哩哔哩直播的请求头
var BiliHeader = map[string]string{
	// referer 必不可少
	"referer": "https://live.bilibili.com/",
	"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/111.0.0.0 Safari/537.36",
}
