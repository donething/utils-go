package dobdpan

// BDFile 网盘的文件。包含上传信息
type BDFile struct {
	Origin       []byte   // 文件的二进制数据
	BlockList    [][]byte // 按每 4MB 分块文件，得到的二进制数据
	BlockMD5List []string // 每个分块的 MD5
	BlockMD5Str  string   // 每个分块的 MD5 数组被转为字符串
	Path         string   // 文件被保存到的远程目录，如"/Pics/filename.jpg"
	Isdir        int
	Size         int
	SliceMd5     string
	ContentMd5   string

	// 可选
	LocalCtime int64 // 创建时间

	// 适配不同网站，手动指定的信息
	Req *Req
}

// Req 百度网盘、一刻相册、Terabox 等不同网站的 API URL 不同，需要手动指定
type Req struct {
	// 上传部分的 URL
	PrecreateURL string
	SuperfileURL string
	CreateURL    string

	ListURL string
	DelURL  string

	// 请求头
	Headers map[string]string
}

//
// 上传文件时的响应
//

// PreResp 响应
type PreResp struct {
	// 不为 0 即表示有错
	Errno int `json:"errno"`

	ReturnType int    `json:"return_type"`
	Uploadid   string `json:"uploadid"`
}

// UpResp 上传分段的响应
type UpResp struct {
	// 为 0 时，在一刻相册中表示有错（只要有 error_code、error_msg 都为有错）；在 Terabox 中非零表示有错
	// 而这些网站都只要 error_msg 不为空""，即表示有错，所以用 error_msg 作为判断标志
	ErrCode int    `json:"error_code"`
	ErrMsg  string `json:"error_msg"`

	Md5 string `json:"md5"`
	// 有时为 string，有时为 int，直接用 Interface{}
	Partseq  interface{} `json:"partseq"`
	Uploadid string      `json:"uploadid"`
}

// CreateResp 创建文件的响应
type CreateResp struct {
	// 不为 0 即表示有错
	Errno int `json:"errno"`

	Data struct {
		Errno          int    `json:"errno"`
		Category       int    `json:"category"`
		FromType       int    `json:"from_type"`
		FSID           int64  `json:"fs_id"`
		Isdir          int    `json:"isdir"`
		Md5            string `json:"md5"`
		Ctime          int64  `json:"ctime"`
		Mtime          int64  `json:"mtime"`
		ShootTime      int64  `json:"shoot_time"`
		Path           string `json:"path"`
		ServerFilename string `json:"server_filename"`
		Size           int64  `json:"size"`
		ServerMd5      string `json:"server_md5"`
	} `json:"data"`
}

// FilesResp 文件列表
type FilesResp struct {
	Cursor string `json:"cursor"`
	Errno  int    `json:"errno"`
	List   []struct {
		Fsid int64 `json:"fsid"`
	} `json:"list"`
}
