package requester

import (
	"io"
	"net/http"
)

// 请求信息
// @author fuzeyu
// @date 2025/2/25
type Client struct {
	httpClient *http.Client
}

// 按照GET请求,获取一个网页信息
// @param uri 请求的uri
func (obj *Client) GetString(uri string) (string, error) {

	resp, err := obj.httpClient.Get(uri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(all), nil
}

func NewClient(httpClient *http.Client) *Client {
	c := &Client{
		httpClient: httpClient,
	}
	return c
}
