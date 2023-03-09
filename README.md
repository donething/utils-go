# utils-go

自用 Golang 工具

# doaudio

发送提示音

```go
err := doaudio.PlayAudio(audios.Beep, 2, 2, BufSize8192)
if err != nil {
  log.Printf("error: %v\n", err)
}
```

# doconf

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

# dodb

## dobadger

使用`badger`键值数据库

```go
const dbDir = "mydb"

var DB *dobadger.DoBadger

func init() {
	// 打开数据库
	DB, err = dobadger.Open(dbDir, nil)
    if err != nil {
        panic(err)
    }
}
```

## dodolt

使用`etcd-io/bbolt`键值数据库

```go
const dbPath = "mydb.db"

var DB *doblot.DoBolt

func init() {
	// 打开数据库
	DB, err = dobolt.Open(dbPath, nil, nil)
    if err != nil {
        panic(err)
    }

	// 创建桶
    err := DB.Create(bucketName)
    if err != nil {
        panic(err)
    }

	// 读取值
    value, err := db.Get([]byte("name"), []byte("people"))
    if err != nil {
        panic(err)
    }
	
	// 写入值
	err := db.Set([]byte("name"), []byte("LiLi"), []byte("people"))
    if err != nil {
        panic(err)
    }
}
```

# dofile

文件操作

# dohttp

执行网络请求

# doint

终端中断命令

# dolog

日志处理

# dotext

文本处理

# dotgpush

TG 消息推送

# dowxpush

微信消息推送

