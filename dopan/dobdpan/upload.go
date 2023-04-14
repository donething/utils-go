// Package dobdpan 上传文件到百度网盘
// @see https://pan.baidu.com/union/document/basic#%E9%A2%84%E4%B8%8A%E4%BC%A0
package dobdpan

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"math/rand"
	"net/url"
	"os"
	"time"
)

const (
	// 按 4MB 分割文件上传
	splitSize = 4 * 1024 * 1024
	// 取前 256 KB  字节计算 MD5
	md5Size = 256 * 1024
)

const (
	tagOneSeq  = `["5910a591dd8fc18c32a8f3df4fdc1761"]`
	tagMoreSeq = `["5910a591dd8fc18c32a8f3df4fdc1761","a5fc157d78e6ad1c7e114b056c92821e"]`
)

var client = dohttp.New(false, false)

// GetYikeReq 获取适配的网站的 Req，将用于创建 BDFile
//
// 参数 bdstoken 一刻相册、百度网盘 需要传递
func GetYikeReq(cookie string, bdstoken string) *Req {
	return &Req{
		PrecreateURL: fmt.Sprintf("https://photo.baidu.com/youai/file/v1/precreate?"+
			"clienttype=70&bdstoken=%s", bdstoken),
		SuperfileURL: "https://c3.pcs.baidu.com/rest/2.0/pcs/superfile2?method=upload&app_id=16051585" +
			"&channel=chunlei&clienttype=70&web=1&logid=MTYyNDAwODkyNzY1NTAuNzEyMjQyOTExODk0OTE1" +
			"&path=%s&uploadid=%s&partseq=%d",
		CreateURL: fmt.Sprintf("https://photo.baidu.com/youai/file/v1/create?"+
			"clienttype=70&bdstoken=%s", bdstoken),
		ListURL: fmt.Sprintf("https://photo.baidu.com/youai/file/v1/list?clienttype=70&"+
			"need_thumbnail=1&need_filter_hidden=0&bdstoken=%s", bdstoken),
		DelURL: fmt.Sprintf("https://photo.baidu.com/youai/file/v1/delete?clienttype=70&"+
			"bdstoken=%s&fsid_list=%s", bdstoken, "%s"),

		Headers: map[string]string{
			"UserAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) " +
				"Chrome/111.0.0.0 Safari/537.36",
			"Origin":  "https://photo.baidu.com",
			"Referer": "https://photo.baidu.com/photo/web/home",
			"Cookie":  cookie,
		},
	}
}

// GetTeraboxReq 获取适配的网站的 Req，将用于创建 BDFile
//
// 参数 bdstoken Terabox 不需要，传到空""即可
func GetTeraboxReq(cookie string) *Req {
	return &Req{
		PrecreateURL: "https://www.terabox.com/api/precreate",
		SuperfileURL: "https://c-jp.terabox.com/rest/2.0/pcs/superfile2?method=upload&app_id=250528&" +
			"channel=dubox&clienttype=0&web=1&logid=MTY3ODc5NjA3MDg0MjAuODU3Mjc0MjM3NzAxNTQ2OA==&" +
			"uploadsign=0&path=%s&uploadid=%s&partseq=%d",
		CreateURL: "https://www.terabox.com/api/create?isdir=0&rtype=1&app_id=250528&web=1&" +
			"channel=dubox&clienttype=0",
		ListURL: "",
		DelURL:  "",

		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) " +
				"Chrome/111.0.0.0 Safari/537.36",
			"Origin":  "https://www.terabox.com",
			"Referer": "https://www.terabox.com",
			"Cookie":  cookie,
		},
	}
}

// NewBytes 根据文件二进制数据，创建 BDFile 的实例
//
// 参数 remotePath 该文件将被保存到的远程路径（可用'/'表示子文件夹）。不过一刻相册中不会区分文件夹
//
// 参数 req 发送请求，不同网站需要的信息。可通过 GetBytes*Req() 快速获取指定网站的 Req 参数
//
// 参数 createdTime 文件被创建的 Unix 时间戳（秒）。为 0 时，将自动设为当前时间戳
func NewBytes(bs []byte, remotePath string, req *Req, createdTime int64) *BDFile {
	// 其它属性
	sec := createdTime
	if sec == 0 {
		sec = time.Now().Unix()
	}

	// 文件对象，将返回
	bdFile := BDFile{
		RemotePath: remotePath,
		LocalCtime: sec,
		Size:       int64(len(bs)),

		Reader: bytes.NewReader(bs),
		Req:    req,
	}

	return &bdFile
}

// NewPath 根据本地文件路径，创建 BDFile 的实例
func NewPath(path string, remotePath string, req *Req) (*BDFile, error) {
	// 提取文件信息
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// 读取文件，循环发送文件的所有切片
	// 注意在 superfile 完后，需要断言为 *File 后调用 close() 关闭 该Reader
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// 文件对象，将返回
	bdFile := BDFile{
		RemotePath: remotePath,
		LocalCtime: 0,
		Size:       fi.Size(),

		Reader: file,
		Req:    req,
	}

	return &bdFile, nil
}

// UploadFile 上传文件到一刻相册
//
// 读取大文件 https://learnku.com/articles/23559/two-schemes-for-reading-golang-super-large-files
func (f *BDFile) UploadFile() error {
	// 关闭由 os.Open() 打开的文件
	defer func() {
		if r, ok := f.Reader.(*os.File); ok {
			r.Close()
		}
	}()

	// 预创建
	resp, err := f.precreate()
	if err != nil {
		return err
	}
	// type 为 2或3，都表示云端已有该文件，经过“预创建”，已经“秒传”，直接返回
	if resp.ReturnType == 2 || resp.ReturnType == 3 {
		return nil
	}
	// 此后，为 1 表示云端没有改文件，需要发送；其它 type 为未知的响应
	if resp.ReturnType != 1 {
		return fmt.Errorf("未知的响应 ReturnType：%+v", resp)
	}

	// 云端没有该文件，需要发送
	// 循环发送文件的所有切片

	// 下面会填满数据，不必 make([]byte, 0, splitSize)
	bs := make([]byte, splitSize)
	seq := 0
	for {
		// 先清空原内容
		n, err := f.Reader.Read(bs[:])
		// 读取出错
		if n < 0 {
			return err
		}
		// 已读完
		if n == 0 {
			break
		}

		// 发送切片
		// 不能直接写 bs，因为最后一次读取不一定填满 bs
		err = f.superfile(bs[0:n], seq, resp.Uploadid)
		if err != nil {
			return err
		}

		f.BlockMD5List = append(f.BlockMD5List, fmt.Sprintf("%x", md5.Sum(bs[0:n])))

		// 继续读取文件
		seq++
	}

	// 已发送所有切片，开始创建文件
	return f.create(resp.Uploadid)
}

// 1. 预处理数据文件
func (f *BDFile) precreate() (*PreResp, error) {
	// 设置该文件是单个切片，还是有多个切片
	if f.Size <= splitSize {
		f.BlockListMd5 = tagOneSeq
	} else {
		f.BlockListMd5 = tagMoreSeq
	}

	// 创建表单
	// "rtype"的值需要为"3"（覆盖文件）
	form := url.Values{}
	form.Add("autoinit", "1")
	form.Add("path", url.QueryEscape(f.RemotePath))

	form.Add("block_list", f.BlockListMd5)
	form.Add("local_ctime", fmt.Sprintf("%d", f.LocalCtime))
	// form.Add("local_mtime", fmt.Sprintf("%d", time.Now().Unix()))

	// 发送表单
	bs, err := client.PostForm(f.Req.PrecreateURL, form.Encode(), f.Req.Headers)
	if err != nil {
		return nil, err
	}

	// 解析
	var resp PreResp
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return &resp, fmt.Errorf("解析 precreate 的响应出错：%w ==> %s", err, string(bs))
	}

	// 响应不符合要求
	if resp.Errno != 0 {
		return &resp, fmt.Errorf("预上传切片失败：%s", string(bs))
	}

	return &resp, nil
}

// 2. 上传切片
//
// @see https://stackoverflow.com/questions/52696921/reading-bytes-into-go-buffer-with-a-fixed-stride-size
func (f *BDFile) superfile(bs []byte, seq int, uploadid string) error {
	// 上传切片
	u := fmt.Sprintf(f.Req.SuperfileURL, url.QueryEscape(f.RemotePath), uploadid, seq)
	file := map[string]interface{}{"file": bs}
	bs, err := client.PostFiles(u, file, nil, f.Req.Headers)
	if err != nil {
		return fmt.Errorf("上传切片出错：%s", err)
	}

	// 解析结果
	var upResp UpResp
	err = json.Unmarshal(bs, &upResp)
	if err != nil {
		return fmt.Errorf("解析上传切片的响应出错：%s", string(bs))
	}

	// 该响应不符合要求
	if upResp.ErrMsg != "" {
		return fmt.Errorf("上传切片失败：%s", string(bs))
	}

	return nil
}

// 3. 根据上传的切片，生成文件
func (f *BDFile) create(uploadid string) error {
	bsMd5List, err := json.Marshal(f.BlockMD5List)
	if err != nil {
		return err
	}

	// 创建表单
	form := url.Values{}
	form.Add("isdir", "0")
	// rtype 的值：1 为重命名同目录、同名文件；3 为始终覆盖
	form.Add("rtype", "1")
	// 不能 url.QueryEscape()，会作为文件名，导致指定的子文件夹失效
	form.Add("path", f.RemotePath)
	form.Add("uploadid", uploadid)
	form.Add("block_list", string(bsMd5List))
	form.Add("size", fmt.Sprintf("%d", f.Size))

	bs, err := client.PostForm(f.Req.CreateURL, form.Encode(), f.Req.Headers)
	if err != nil {
		return fmt.Errorf("创建文件出错：%s", err)
	}

	var cResp CreateResp
	err = json.Unmarshal(bs, &cResp)
	if err != nil {
		return fmt.Errorf("解析创建文件的响应出错：%s", string(bs))
	}

	if cResp.Errno != 0 {
		return fmt.Errorf("创建文件失败：%s", string(bs))
	}

	return nil
}

// DelAll 删除所有文件
func DelAll(req *Req) error {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		// 列出文件
		bs, err := client.GetBytes(req.ListURL, req.Headers)
		if err != nil {
			return fmt.Errorf("列出文件出错：%w", err)
		}

		var files FilesResp
		err = json.Unmarshal(bs, &files)
		if err != nil {
			return fmt.Errorf("解析文件列表出错：%w", err)
		}

		if files.Errno != 0 {
			return fmt.Errorf("列出文件失败：%s\n", string(bs))
		}

		// 删除文件
		fidList := make([]int64, len(files.List))
		for i, f := range files.List {
			fidList[i] = f.Fsid
		}
		bs, err = json.Marshal(fidList)
		if err != nil {
			return fmt.Errorf("序列化文件的 ID 列表时出错：%w", err)
		}

		u := fmt.Sprintf(req.DelURL, string(bs))
		bs, err = client.GetBytes(u, req.Headers)
		if err != nil {
			return fmt.Errorf("删除文件出错：%w", err)
		}

		var resp PreResp
		err = json.Unmarshal(bs, &resp)
		if err != nil {
			return fmt.Errorf("解析删除文件的响应时出错：%w", err)
		}

		if resp.Errno == 2 {
			break
		}
		if resp.Errno != 0 {
			return fmt.Errorf("删除一刻相册中所有的图片失败：%s", string(bs))
		}

		fmt.Printf("已删除该页图片，将继续删除下页\n")

		r := rand.Intn(5)
		time.Sleep(time.Duration(r+1) * time.Second)
	}

	return nil
}
