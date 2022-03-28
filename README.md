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

