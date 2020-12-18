package requests

import (
	"encoding/json"
	"net/http"
	"time"
)

//负责response
type Response struct {
	client     *Client        //原始request 对象
	request    *http.Request  //httprequest对象
	resp       *http.Response //http.Response对象
	httpClient *http.Client   //client对象
	timeCost   time.Duration
	reqBody    []byte //请求body
	respBody   []byte //返回body
	Err        error  //出错信息
}

//将response body解析到 v的结构体中
func (resp *Response) GetJson(v interface{}) error {
	data := resp.respBody
	return json.Unmarshal(data, v)
}

func (resp *Response) GetBodyBytes() []byte {
	return resp.respBody
}

func (resp *Response) GetTimeCost() time.Duration {
	return resp.timeCost
}

func (resp *Response) GetBodyString() string {
	return string(resp.respBody)
}

//如果不满足需求  可获取原生httpResponse
func (resp *Response) GetHttpResponse() *http.Response {
	return resp.resp
}

func (resp *Response) GetHeader(key string) []string{
	headerVal := resp.resp.Header[key]
	return headerVal

}
