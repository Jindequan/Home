package requests

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"
)

//对外暴露http库的细节 方便requests库中考虑不全面的外面可以直接拿client进行操作

//返回默认参数的client
func GetDefaultHttpClient() *http.Client {
	jar, _ := cookiejar.New(nil)
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &http.Client{
		Jar:       jar,
		Transport: transport,
		Timeout:   2 * time.Minute,
	}
}

func (r *Client) SetClient(client *http.Client) {
	r.client = client // use default if client == nil
}

func (r *Client) GetClient() *http.Client {
	return r.client
}
