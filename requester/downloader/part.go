package downloader

// 分片信息
type part struct {

	// 起始和结束范围
	startOffset int64
	endOffset   int64

	// 当前是第几片
	partNumber int

	// 下载内容
	buffers []byte

	// 状态
	status bufferStatus
}
