package requester

import (
	"github.com/chaksunshine/kit/json"
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

// 获取get请求的buffer对象
// @param uri 请求的uri
func (obj *Client) getBuffers(uri string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(thread.CtxRequest(), "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	obj.formatRequester(request)

	do, err := obj.httpClient.Do(request)
	return do, err
}

// 按照GET请求,获取一个网页信息
// @param uri 请求的uri
func (obj *Client) GET(uri string) ([]byte, error) {
	buffers, err := obj.getBuffers(uri)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(buffers.Body)
}

// 按照GET请求,获取一个网页信息,并格式化为json
// @param uri 请求的uri
// @param data 数据格式
// @param getBuffers 获取buffer内容
func (obj *Client) GETFormtJson(uri string, data interface{}, loadBuffers ...bool) ([]byte, error) {

	buffers, err := obj.getBuffers(uri)
	if err != nil {
		return nil, err
	}

	// 获取buffer内容
	if len(loadBuffers) == 1 && loadBuffers[0] {

		// 获取内容
		all, err := io.ReadAll(buffers.Body)
		if err != nil {
			return nil, err
		}

		// 解析内容
		if err := json.Unmarshal(all, &data); err != nil {
			return nil, err
		}
		return all, nil
	}

	// 不获取buffer内容
	decoder := json.NewDecoder(buffers.Body)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return nil, nil
}

// @param httpClient 客户端信息
func NewClient(httpClient *http.Client) *Client {
	c := &Client{
		httpClient: httpClient,
		config:     defaultConfig,
	}
	return c
}
