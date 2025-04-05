package downloader

import "errors"

// 检查是否支持分片下载
var errNotAllowBytes = errors.New("不支持分片下载")

// 每片文件大小
type PartSize int

const (

	// 5M一片
	PartSizeSmall = 1024 * 1024 * 5
)

// 下载状态
type bufferStatus int

const (

	// 等待被下载
	bufferStatusWaitWaitDownload bufferStatus = iota

	// 下载完成
	bufferStatusWaitWriteFile

	// 需要重试
	bufferStatusRetry
)
