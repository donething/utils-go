// Package dobdpan 上传文件到百度网盘
// @see https://pan.baidu.com/union/document/basic#%E9%A2%84%E4%B8%8A%E4%BC%A0
package dobdpan

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"math"
	"math/rand"
	"net/url"
	"time"
)

const (
	// 按 4MB 分割文件上传
	splitSize = 4 * 1024 * 1024
	// 取前 256 KB  字节计算 MD5
	md5Size = 256 * 1024
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
		DelURL: "https://photo.baidu.com/youai/file/v1/delete?clienttype=70&bdstoken=%s&fsid_list=%s",

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

// New 创建 BDFile 的实例
//
// 参数 data 要上传文件的二进制数据
//
// 参数 remotePath 该文件将被保存到的远程路径（可用'/'表示子文件夹）。不过一刻相册中不会区分文件夹
//
// 参数 createdTime 文件被创建的 Unix 时间戳（秒）。为 0 时，将自动设为当前时间戳
//
// 参数 req 发送请求，不同网站需要的信息。可通过 Get*Req() 快速获取指定网站的 Req 参数
func New(data []byte, remotePath string, createdTime int64, req *Req) *BDFile {
	// 文件将被分成的段数
	blockNum := int(math.Ceil(float64(len(data)) / float64(splitSize)))

	// 当文件大小小于 md5Size 时，两个 MD5 相同
	var contentMd5 = md5.Sum(data)
	var sliceMd5 = contentMd5
	if len(data) > md5Size {
		sliceMd5 = md5.Sum(data[:md5Size])
	}

	// 文件对象，将返回
	bdFile := BDFile{
		Origin:       data,
		BlockList:    make([][]byte, blockNum),
		BlockMD5List: make([]string, blockNum),
		Isdir:        0,
		Path:         remotePath,
		Size:         len(data),
		ContentMd5:   fmt.Sprintf("%x", contentMd5),
		SliceMd5:     fmt.Sprintf("%x", sliceMd5),

		Req: req,
	}

	// 其它属性
	sec := createdTime
	if sec == 0 {
		sec = time.Now().Unix()
	}
	bdFile.LocalCtime = sec

	// 将文件分段
	i := 0
	for pos := 0; i < blockNum; pos += splitSize {
		var tmp []byte
		// 前面的分段为 [pos : pos+splitSize]，最后一个分段为 [pos:]
		if i <= blockNum-2 {
			tmp = data[pos : pos+splitSize]
		} else {
			tmp = data[pos:]
		}

		// 添加分段
		bdFile.BlockList[i] = tmp
		// 保存 MD5
		bdFile.BlockMD5List[i] = fmt.Sprintf("%x", md5.Sum(tmp))
		i++
	}

	// 将 md5 的数组转为字符串
	md5BS, _ := json.Marshal(bdFile.BlockMD5List)
	bdFile.BlockMD5Str = string(md5BS)

	return &bdFile
}

// UploadFile 上传文件到一刻相册
func (f *BDFile) UploadFile() error {
	resp, err := f.precreate()
	if err != nil {
		return err
	}

	// type 为 1，表示云端没有该文件，需要上传
	if resp.ReturnType == 1 {
		// 上传
		err = f.superfile(resp)
		if err != nil {
			return err
		}

		err = f.create(resp.Uploadid)
		return err
	} else if resp.ReturnType == 2 || resp.ReturnType == 3 {
		// type 为 2或3，都表示云端已有该文件，可以“秒传”
		return nil
	}

	return fmt.Errorf("未知的响应 ReturnType：%+v", resp)
}

// 1. 预处理数据文件
func (f *BDFile) precreate() (*PreResp, error) {
	// 创建表单
	// "rtype"的值需要为"3"（覆盖文件）
	form := url.Values{}
	form.Add("autoinit", "1")
	form.Add("isdir", fmt.Sprintf("%d", f.Isdir))
	form.Add("rtype", "3")
	form.Add("ctype", "11")
	form.Add("path", f.Path)

	form.Add("content-md5", f.ContentMd5)
	form.Add("size", fmt.Sprintf("%d", f.Size))
	form.Add("slice-md5", f.SliceMd5)
	form.Add("block_list", f.BlockMD5Str)
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
		return &resp, fmt.Errorf("预上传分段失败：%s", string(bs))
	}

	return &resp, nil
}

// 2. 分段上传
//
// @see https://stackoverflow.com/questions/52696921/reading-bytes-into-go-buffer-with-a-fixed-stride-size
func (f *BDFile) superfile(resp *PreResp) error {
	for i := 0; i < len(f.BlockList); i++ {
		// 上传片段
		u := fmt.Sprintf(f.Req.SuperfileURL, url.QueryEscape(f.Path), resp.Uploadid, i)
		file := map[string]interface{}{"file": f.BlockList[i]}
		bs, err := client.PostFiles(u, file, nil, f.Req.Headers)
		if err != nil {
			return fmt.Errorf("上传分段出错：%w ==> %s", err, string(bs))
		}

		// 解析结果
		var upResp UpResp
		err = json.Unmarshal(bs, &upResp)
		if err != nil {
			return fmt.Errorf("解析上传分段的响应出错：%w ==> %s", err, string(bs))
		}

		// 该响应不符合要求
		if upResp.ErrMsg != "" {
			return fmt.Errorf("上传分段失败：%s", string(bs))
		}
	}

	return nil
}

// 3. 根据上传的分段，生成文件
func (f *BDFile) create(uploadid string) error {
	// 创建表单
	form := url.Values{}
	form.Add("isdir", fmt.Sprintf("%d", f.Isdir))
	form.Add("rtype", "3")
	form.Add("ctype", "11")
	form.Add("path", f.Path)

	form.Add("content-md5", f.ContentMd5)
	form.Add("size", fmt.Sprintf("%d", f.Size))
	form.Add("uploadid", uploadid)
	form.Add("block_list", f.BlockMD5Str)

	bs, err := client.PostForm(f.Req.CreateURL, form.Encode(), f.Req.Headers)
	if err != nil {
		return fmt.Errorf("创建文件出错：%w ==> %s", err, string(bs))
	}

	var cResp CreateResp
	err = json.Unmarshal(bs, &cResp)
	if err != nil {
		return fmt.Errorf("解析创建文件的响应出错：%w ==> %s", err, string(bs))
	}

	if cResp.Errno != 0 {
		return fmt.Errorf("创建文件失败：%s", string(bs))
	}

	return nil
}

// DelAll 删除所有文件
func DelAll(req *Req, bdstoken string) error {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		// 列出文件
		bs, err := client.Get(req.ListURL, req.Headers)
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

		u := fmt.Sprintf(req.DelURL, bdstoken, string(bs))
		bs, err = client.Get(u, req.Headers)
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
