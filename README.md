# Go YouTube Downloader (High Quality)

這是一個使用 Golang 開發的高畫質 YouTube 影片下載工具。它能自動抓取最高畫質的視訊與最高音質的音訊，並透過 FFmpeg 進行合併。

## ✨ 特色
- **4K/1080p 支援**：自動篩選最高畫質。
- **併發下載**：使用 Goroutine 同時下載影音，速度提升一倍。
- **進度條顯示**：即時掌握下載百分比與速度。
- **安全檔名**：自動過濾 Windows 不允許的特殊字元。

## 🛠️ 環境需求
1. **Golang**: 1.18+
2. **FFmpeg**: 必須安裝於系統中（本程式會調用 `ffmpeg` 指令）。

## 🚀 快速開始

### 1. 安裝套件
```bash
go mod tidy

### 2. 執行下載
go run main.go -url "你的 YouTube 網址"

### 3. 編譯為執行檔
go build -o ytdl.exe main.go
.\ytdl.exe -url "你的 YouTube 網址"

📝 注意事項
請確保 ffmpeg.exe 的路徑在程式碼中設定正確（預設為 C:\ffmpeg\bin\ffmpeg.exe）。

本工具僅供學習與技術交流使用，請遵守 YouTube 相關服務條款。