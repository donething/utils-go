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
)

// Cut 切割视频为多个分段
//
// maxSegSize 单位字节
//
// 当 dstPath 目标路径为空""时，默认保存到视频的同目录下
func Cut(path string, maxSegSize int64, dstDir string) error {
	// 默认保存到视频的同目录下
	if strings.TrimSpace(dstDir) == "" {
		dstDir = filepath.Dir(path)
	}

	err := os.MkdirAll(dstDir, 0755)
	if err != nil {
		return err
	}

	file, err := os.Stat(path)
	if err != nil {
		return err
	}

	// 计算分段数
	n := int(math.Ceil(float64(file.Size()) / float64(maxSegSize)))

	seconds, err := GetDuration(path)
	if err != nil {
		return err
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
		return fmt.Errorf("%w: %s", err, string(output))
	}

	return nil
}

// Convt 转换编码
//
// 如果 dstPath 目标路径为空""，将转码为".mp4"，并保存到视频同目录下
func Convt(path string, dstPath string) error {
	// 默认转码为".mp4"，并保存到视频同目录下
	if strings.TrimSpace(dstPath) == "" {
		dstPath = strings.TrimSuffix(path, filepath.Ext(path)) + ".mp4"
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
		return fmt.Errorf("%w: %s", err, string(output))
	}

	return nil
}

// GetDuration 获取视频的时长（秒）
func GetDuration(path string) (int, error) {
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
		return 0, fmt.Errorf("%w: %s", err, string(output))
	}

	str := strings.TrimSpace(string(output))
	seconds, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}

	return int(math.Ceil(seconds)), nil
}

// GetResolution 获取视频的分辨率，返回：宽度、高度
func GetResolution(path string) (width int, height int, err error) {
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
		return 0, 0, fmt.Errorf("%w: %s", err, string(output))
	}

	dimensions := strings.Split(strings.TrimSpace(string(output)), "x")
	width, err = strconv.Atoi(dimensions[0])
	if err != nil {
		return
	}
	height, err = strconv.Atoi(dimensions[1])
	if err != nil {
		return
	}

	return
}

// GetFrame 获取指定时刻的帧。推荐保存为 .jpg 文件
//
// time 时间戳，如"01:20:10"
//
// resolution 宽高比，如"640:480"
func GetFrame(path string, dstPath string, time string, resolution string) error {
	args := []string{
		"-hide_banner",
		"-i", path,
		"-ss", time,
		"-vframes", "1",
		"-vf", fmt.Sprintf("scale=%s:force_original_aspect_ratio=decrease", resolution),
		dstPath,
	}

	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w:%s", err, string(output))
	}

	return nil
}
