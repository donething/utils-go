# utils-go

自用 Golang 工具

## doaudio

发送提示音

```go
err := doaudio.PlayAudio(audios.Beep, 2, 2, BufSize8192)
if err != nil {
  log.Printf("error: %v\n", err)
}
```

## doconf

```go
type DoConf struct {
Name  string `json:"name"`
Age   int    `json:"age"`
Extra string `json:"extra"`
}

const confPath = "./test_conf.json"

var Conf DoConf

func initConf() {
    exist, err := doconf.Init(confPath, &Conf)
    if err != nil {
        t.Fatal(err)
    }

    if !exist {
        t.Logf("似乎是首次运行，请先填写配置后，再运行程序")
        return
    }

    t.Logf("配置：%+v\n", Conf)
}
```

## dodb

doltdb 键值数据库操作

## dofile

文件操作

## dohttp

执行网络请求

## doint

终端中断命令

## dolog

日志处理

## dotext

文本处理

## dotgpush

TG 消息推送

## dowxpush

微信消息推送

