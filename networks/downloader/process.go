package downloader

import (
	"context"
	"errors"
	"github.com/chaksunshine/kit/numbers"
	"net/http"
	"os"
)

// 分片下载起
// @author fuzeyu
// @date 2025/2/23
type Process struct {
	ctx context.Context
	uri string

	// 分片大小
	partSize PartSize
	partList []*part
}

// 将一个总大小的文件分割成多片
// @param totalSize 文件大小
func (obj *Process) parserBuffers(totalSize int64) {

	var partIndex = 1

	chunkSize := int64(obj.partSize)
	for start := int64(0); start < totalSize; start += chunkSize {
		end := start + chunkSize - 1
		if end > totalSize {
			end = totalSize // 最后一个分片可能不足 chunkSize
		}

		partInfo := &part{
			startOffset: start,
			endOffset:   end,
			partNumber:  partIndex,
			buffers:     make([]byte, 0),
			status:      bufferStatusWaitWaitDownload,
		}
		obj.partList = append(obj.partList, partInfo)
		partIndex++
	}
}

// 将任务分片下载
func (obj *Process) parserPart() error {

	resp, err := http.DefaultClient.Head(obj.uri)
	if err != nil {
		return err
	}
	header := resp.Header

	// 检查是否支持分片下载
	ranges := header.Get("Accept-Ranges")
	if ranges != "bytes" {
		return errNotAllowBytes
	}

	// 检查包大小
	contentLength := header.Get("Content-Length")
	contentSize := numbers.StringToInt64Must(contentLength)
	if contentSize <= 0 {
		return errors.New("未获取到包的大小")
	}
	obj.parserBuffers(contentSize)

	return nil
}

// 执行
func (obj *Process) execute() error {

	// 解析分片信息
	if err := obj.parserPart(); err != nil {
		return err
	}

	// 开始下载

	os.Exit(2)
	return nil
}

// @param ctx 上下文
// @param uri 下载地址
// @param partSize 分片大小
func NewProcess(ctx context.Context, uri string, partSize PartSize) (*Process, error) {
	c := &Process{
		ctx:      ctx,
		uri:      uri,
		partSize: partSize,
	}
	return c, c.execute()
}
