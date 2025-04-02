package requester

import (
	"github.com/chaksunshine/kit/thread"
	"io"
	"net/http"
)

// 请求信息
// @author fuzeyu
// @date 2025/2/25
type Client struct {
	httpClient *http.Client
	config     *Config
}

// 格式化请求信息
// @param request 请求信息
func (obj *Client) formatRequester(request *http.Request) {

	// 添加请求头
	if obj.config.FormatHeader {
		request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
		request.Header.Add("Cache-Control", "max-age=0")
		request.Header.Add("Dnt", "1")
		request.Header.Add("Sec-Ch-Ua", `"Microsoft Edge";v="125", "Chromium";v="125", "Not.A/Brand";v="24"`)
		request.Header.Add("Sec-Ch-Ua-Mobile", "?0")
		request.Header.Add("Sec-Ch-Ua-Platform", `"Windows"`)
		request.Header.Add("Sec-Fetch-Dest", "document")
		request.Header.Add("Sec-Fetch-Mode", "navigate")
		request.Header.Add("Sec-Fetch-Site", "none")
		request.Header.Add("Sec-Fetch-User", "?1")
		request.Header.Add("Upgrade-Insecure-Requests", "1")
		request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0")
	}
}

// 按照GET请求,获取一个网页信息
// @param uri 请求的uri
func (obj *Client) GetString(uri string) (string, error) {

	request, err := http.NewRequestWithContext(thread.CtxRequest(), "GET", uri, nil)
	if err != nil {
		return "", err
	}
	obj.formatRequester(request)

	do, err := obj.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer do.Body.Close()

	all, err := io.ReadAll(do.Body)
	if err != nil {
		return "", err
	}
	return string(all), nil
}

// @param httpClient 客户端信息
func NewClient(httpClient *http.Client) *Client {
	c := &Client{
		httpClient: httpClient,
		config:     defaultConfig,
	}
	return c
}
