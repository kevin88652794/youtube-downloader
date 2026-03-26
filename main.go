package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/kkdai/youtube/v2"
	"github.com/schollz/progressbar/v3"
)

func main() {
	videoURL := flag.String("url", "", "YouTube 影片網址")
	flag.Parse()

	if *videoURL == "" {
		fmt.Println("❌ 請提供網址: -url=https://...")
		return
	}

	client := youtube.Client{}
	video, err := client.GetVideo(*videoURL)
	if err != nil {
		fmt.Printf("❌ 取得影片資訊失敗: %v\n", err)
		return
	}

	// 1. 安全檔名處理
	re := regexp.MustCompile(`[\\/:*?"<>|]`)
	safeTitle := re.ReplaceAllString(video.Title, "")
	fmt.Printf("🎬 準備下載: %s\n", safeTitle)

	tempDir := "./temp_cache"
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	vFile := filepath.Join(tempDir, "v.mp4")
	aFile := filepath.Join(tempDir, "a.mp4")
	finalFile := safeTitle + ".mp4"

	// 2. 篩選最佳格式
	var bestV, bestA *youtube.Format
	for i := range video.Formats {
		f := &video.Formats[i]
		if f.AudioChannels == 0 && f.QualityLabel != "" {
			if bestV == nil || f.Bitrate > bestV.Bitrate {
				bestV = f
			}
		}
		if f.QualityLabel == "" {
			if bestA == nil || f.Bitrate > bestA.Bitrate {
				bestA = f
			}
		}
	}

	if bestV == nil || bestA == nil {
		fmt.Println("❌ 找不到適當的影音格式")
		return
	}

	// 3. 併發下載與進度條
	var wg sync.WaitGroup
	wg.Add(2)
	errChan := make(chan error, 2)

	go func() {
		defer wg.Done()
		if err := downloadWithProgress(client, video, bestV, vFile, "📹 視訊"); err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := downloadWithProgress(client, video, bestA, aFile, "🎵 音訊"); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for e := range errChan {
		if e != nil {
			fmt.Printf("\n❌ 下載失敗: %v\n", e)
			return
		}
	}

	// 4. FFmpeg 合併 (直接指定絕對路徑)
	fmt.Println("\n⚙️  FFmpeg 合併影音軌中...")

	// 請確保 ffmpeg.exe 確實位於此路徑
	ffmpegPath := `C:\ffmpeg\bin\ffmpeg.exe`

	cmd := exec.Command(ffmpegPath, "-y", "-i", vFile, "-i", aFile, "-c", "copy", finalFile)

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ 合併失敗，請確認 FFmpeg 路徑正確: %v\n", err)
		return
	}

	fmt.Printf("\n✨ 下載完成！檔案存於: %s\n", finalFile)
}

func downloadWithProgress(client youtube.Client, video *youtube.Video, format *youtube.Format, path string, desc string) error {
	stream, size, err := client.GetStream(video, format)
	if err != nil {
		return err
	}
	defer stream.Close()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	bar := progressbar.DefaultBytes(size, desc)
	_, err = io.Copy(io.MultiWriter(file, bar), stream)
	return err
}
