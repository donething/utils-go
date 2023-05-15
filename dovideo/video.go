// Package dovideo 使用 FFmpeg 处理视频（务必安装了 FFmpeg，并添加了系统环境变量）
package dovideo

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNoVideos = fmt.Errorf("不存在需要合并的视频片段")
)

// Cut 切割视频为多个分段
//
// maxSegSize 单位字节
//
// 当 dstPath 目标路径为空""时，默认保存到视频的同目录下
//
// 注意：因为 ffmpeg 切割视频的时间点并不准确，切割出来的文件数量，不一定等于`文件字节数/分段大小`，
// 所以返回的路径列表中，有的路径不存在。所以需要专门判断，去除不存在的路径
func Cut(path string, maxSegSize int64, dstDir string) ([]string, error) {
	tag := "Cut"
	// 默认保存到视频的同目录下
	if strings.TrimSpace(dstDir) == "" {
		dstDir = filepath.Dir(path)
	}

	err := os.MkdirAll(dstDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("[%s]创建临时目录出错：%w", tag, err)
	}

	file, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("[%s]获取待切割视频的文件信息出错：%w", tag, err)
	}

	// 计算分段数
	n := int(math.Ceil(float64(file.Size()) / float64(maxSegSize)))

	seconds, err := GetDuration(path)
	if err != nil {
		return nil, fmt.Errorf("[%s]获取视频时长出错：%w", tag, err)
	}

	// 每个分段的时长（秒）
	segmentDuration := seconds / n

	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	args := []string{
		"-hide_banner",
		"-i", path,
		"-c", "copy",
		"-f", "segment",
		"-segment_time", fmt.Sprintf("%d", segmentDuration),
		// 需要重置时间戳，否则每个片段的进度条依然是原视频的长度
		"-reset_timestamps", "1",
		// 从 1 开始编号，而不是默认的 0
		"-segment_start_number", "1",
		filepath.Join(dstDir, name+"_%02d.mp4"),
	}
	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("[%s]执行切割视频出错：%w: %s", tag, err, string(output))
	}

	// 用于返回路径的数组
	dstPaths := make([]string, 0)
	for i := 0; i < n; i++ {
		segPath := filepath.Join(dstDir, fmt.Sprintf("%s_%02d.mp4", name, i+1))
		// 可以参考函数的使用说明，因为切片不严格按照时间，所以后面的路径可能不存在
		// 也因为这个原因，最后一个分段可能极小（80KB），几乎不含数据，应该忽略，也避免提取缩略图时失败
		if fi, err := os.Stat(segPath); os.IsNotExist(err) || fi.Size() < 512*1024 {
			fmt.Printf("忽略极小的视频分段(%d): %s\n", i+1, segPath)
			// 移除可能极小分段
			_ = os.Remove(segPath)
			break
		}

		dstPaths = append(dstPaths, segPath)
	}

	return dstPaths, nil
}

// Concat 合并视频
//
// dir 指定视频片段所在的目录
//
// inputFormat 指定只包含该格式的视频片段。如 ".ts"
//
// outputFormat 合并后输入的视频格式。如 ".mp4"
//
// 当目录下没有需要合并的视频片段时，返回错误 `ErrNoVideos`
func Concat(dir string, inputFormat string, outputFormat string) (string, error) {
	tag := "Concat"
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("[%s] 读取目录出错：%w", tag, err)
	}

	// 需要合并的视频文件的路径
	videosPath := ""
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// 只包括指定输入格式的视频文件
		if !strings.HasSuffix(file.Name(), inputFormat) {
			continue
		}

		videosPath += fmt.Sprintf("file '%s'\n", filepath.Join(dir, file.Name()))
	}

	// 没有需要合并的视频片段
	if videosPath == "" {
		return "", ErrNoVideos
	}

	// 写入路径到文件
	listFilePath := filepath.Join(dir, fmt.Sprintf("files_%d.txt", time.Now().UnixMilli()))
	err = os.WriteFile(listFilePath, []byte(videosPath), 0644)
	if err != nil {
		return "", fmt.Errorf("[%s]写入视频的路径出错：%w", tag, err)
	}

	// 	调用 FFmpeg 合并视频
	outputPath := filepath.Join(dir, fmt.Sprintf("output_%d%s", time.Now().UnixMilli(), outputFormat))
	args := []string{
		"-hide_banner",
		"-f", "concat",
		"-safe", "0",
		"-i", listFilePath,
		"-c", "copy",
		outputPath,
	}
	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("[%s]执行合并视频出错：%w: %s", tag, err, string(output))
	}

	return outputPath, nil
}

// Convt 转换编码
//
// 如果 dstPath 目标路径为空""，将转码为".mp4"，并保存到视频同目录下
func Convt(path string, dstPath string) error {
	tag := "Convt"
	// 默认转码为".mp4"，并保存到视频同目录下
	if strings.TrimSpace(dstPath) == "" {
		ext := filepath.Ext(path)
		if strings.ToLower(ext) == ".mp4" {
			return fmt.Errorf("[%s]视频已经是 .mp4，无法按默认转为 .mp4。请指定 dstPath 参数", tag)
		}

		dstPath = strings.TrimSuffix(path, ext) + ".mp4"
	}

	args := []string{
		"-hide_banner",
		"-i", path,
		"-c", "copy",
		dstPath,
	}
	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("[%s]执行转换视频出错'%s'：%w: %s", tag, path, err, string(output))
	}

	return nil
}

// GetDuration 获取视频的时长（秒）
func GetDuration(path string) (int, error) {
	tag := "GetDuration"
	// 构建命令行参数
	args := []string{
		"-i", path,
		"-select_streams", "v:0",
		"-show_entries", "format=duration",
		"-v", "quiet",
		"-of", "csv=p=0",
	}

	// 执行命令行命令
	cmd := exec.Command("ffprobe", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("[%s]执行获取视频时长出错：%w: %s", tag, err, string(output))
	}

	str := strings.TrimSpace(string(output))
	seconds, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, fmt.Errorf("[%s]解析时长出错：%w", tag, err)
	}

	return int(math.Ceil(seconds)), nil
}

// GetResolution 获取视频的分辨率，返回：宽度、高度
func GetResolution(path string) (width int, height int, err error) {
	tag := "GetResolution"
	args := []string{
		"-i", path,
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height",
		"-v", "quiet",
		"-of", "csv=s=x:p=0",
	}

	// 执行命令
	cmd := exec.Command("ffprobe", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, 0, fmt.Errorf("[%s]执行获取视频分辨率出错：%w: %s", tag, err, string(output))
	}

	dimensions := strings.Split(strings.TrimSpace(string(output)), "x")
	width, err = strconv.Atoi(dimensions[0])
	if err != nil {
		return 0, 0, fmt.Errorf("[%s]转换宽度值出错'%s'：%w", tag, dimensions[0], err)
	}
	height, err = strconv.Atoi(dimensions[1])
	if err != nil {
		return 0, 0, fmt.Errorf("[%s]转换高度值出错'%s'：%w", tag, dimensions[1], err)
	}

	return
}

// GetFrame 获取指定时刻的帧。推荐保存为 .jpg 文件
//
// time 时间戳，如"01:20:10"
//
// resolution 缩放的宽高比，如"640:480"。为空""则缩放
func GetFrame(path string, dstPath string, time string, resolution string) error {
	tag := "GetFrame"
	args := []string{
		"-y", "-hide_banner",
		"-i", path,
		"-ss", time,
		"-vframes", "1",
	}
	if resolution != "" {
		r := fmt.Sprintf("scale=%s:force_original_aspect_ratio=decrease", resolution)
		args = append(args, []string{"-vf", r}...)
	}

	args = append(args, dstPath)

	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("[%s]执行获取视频某时刻帧出错：%w:%s", tag, err, string(output))
	}

	return nil
}
