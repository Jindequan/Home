package requests

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//负责处理请求的的go文件

const (
	//debug信息
	HeaderDump   = 1 << iota //是否打印请求头部信息
	BodyDump                 //是否打印请求体信息
	RespHeadDump             //是否打印响应头信息
	RespBodyDump             //是否打印请求体信息
	TimeCostDump             //是否打印请求耗时
)

const AllDebugSign = HeaderDump | BodyDump | RespHeadDump | RespBodyDump | TimeCostDump

var defaultClient *Client

//直接用这两个类型 我会把content-type设置为application/json
type JsonBodyStr string
type JsonBodyBytes []byte

var ErrTimeOut = errors.New("请求超时")
var ErrCanRetry = errors.New("请求错误 但是可以重试")

//拼接在URL中的参数 需要配置时 初始化一个QueryString类型的参数传入 ex:req.Do(method,url,QueryString{"a":1,"b":2})
type QueryString map[string]interface{}

//Host 请求中的Host 需要配置时 初始化一个Host类型的参数 传入 ex var host Host = "www.baidu.com" req.Do(method,url,host)
type Host string

//提交表单数据时  初始化一个FormData的参数  ex:req.Do(method,url,FormData{"a":1,"b":2})
type FormData map[string]interface{}

//请求中的header ex:req.Do(method,url,Header{"a":"1","b":"2"})
type Header map[string]string

type RetryCondition func(*Client, *Response) bool

type RetryTime int //重试次数

type param struct {
	url.Values
}

type FileUpload struct {
	// filename in multipart form.
	FileName string
	// form field name
	FieldName string
	// file to uplaod, required
	File io.Reader
}

type BeforeRequestHook struct {
	handle func(client *Client, response *Response)
}
type AfterRequestHook struct {
	handle func(client *Client, response *Response)
}

//单例
func GetDefaultClient() *Client {
	if defaultClient == nil {
		defaultClient = New(GetDefaultHttpClient())
	}
	return defaultClient
}

func (p *param) getValues() url.Values {
	if p.Values == nil {
		p.Values = make(url.Values)
	}
	return p.Values
}

func (p *param) AddMap(addedMap map[string]interface{}) {
	if len(addedMap) == 0 {
		return
	}
	vs := p.getValues()

	for k, v := range addedMap {
		switch v.(type) {
		case []string, []int:
			for _, addedV := range v.([]string) {
				vs.Add(k, fmt.Sprint(addedV))
			}
		default:
			vs.Add(k, fmt.Sprint(v))
		}
	}
}

func (p *param) Copy(beCopiedParam param) {
	if beCopiedParam.Values == nil {
		return
	}
	oriVals := p.getValues()
	for key, values := range beCopiedParam.Values {
		for _, value := range values {
			oriVals.Add(key, value)
		}
	}
}

func (p *param) Empty() bool {
	return p.Values == nil
}

//type
type Client struct {
	client    *http.Client //http库的client对象 不对外暴露
	debugFlag uint8
}

func (client *Client) SetDebugFlag(flag uint8) {
	client.debugFlag = flag
}

func New(client *http.Client) *Client {
	if client == nil {
		return &Client{client: &http.Client{}}
	}

	return &Client{client: client}

}

//接受[]byte和string和*bytes.Buffer
func setBody(req *http.Request, data interface{}) error {
	var handledData []byte
	switch body := data.(type) {
	case string:
		handledData = []byte(body)
	case []byte:
		handledData = body
	case *bytes.Buffer:
		handledData = body.Bytes()
	default:
		return errors.New("不支持的body类型")
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(handledData))
	req.ContentLength = int64(len(handledData))
	return nil
}

func (client *Client) Do(method, reqUrl string, vals ...interface{}) *Response {
	if reqUrl == "" {
		return &Response{Err: errors.New("url不能为空")}
	}

	method = strings.ToUpper(method)

	req := &http.Request{
		Method: method,
		Header: make(http.Header),
	}

	response := &Response{client: client, request: req, httpClient: client.client}

	var queryString param
	var formParam param
	var upload *FileUpload
	var beforeRequestHooks []func(client *Client, response *Response)
	var afterRequestHooks []func(client *Client, response *Response)
	var retryFunc RetryCondition

	for _, v := range vals {
		switch content := v.(type) {
		case Header:
			for headerKey, headerVal := range content {
				req.Header.Add(headerKey, headerVal)
			}
		case http.Header:
			for headerKey, HeaderVals := range content {
				for _, headerVal := range HeaderVals {
					req.Header.Add(headerKey, headerVal)
				}
			}
		case QueryString:
			queryString.AddMap(content)
		case FormData:
			formParam.AddMap(content)
		case string:
			err := setBody(req, content)
			if err != nil {
				response.Err = err
				return response
			}
		case []byte:
			err := setBody(req, content)
			if err != nil {
				response.Err = err
				return response
			}
		case *bytes.Buffer:
			err := setBody(req, content)
			if err != nil {
				response.Err = err
				return response
			}
		case RetryCondition:
			retryFunc = content
		case JsonBodyStr:
			err := setBody(req, []byte(content))
			if err != nil {
				response.Err = err
				return response
			}
			setContentType(req, "application/json")
		case *FileUpload:
			upload = content
		case JsonBodyBytes:
			err := setBody(req, []byte(content))
			if err != nil {
				response.Err = err
				return response
			}
			setContentType(req, "application/json")

		case *http.Cookie:
			req.AddCookie(content)
		case context.Context:
			req.WithContext(content)
		case *BeforeRequestHook:
			beforeRequestHooks = append(beforeRequestHooks, content.handle)
		case *AfterRequestHook:
			afterRequestHooks = append(afterRequestHooks, content.handle)
		case Host:
			req.Host = string(content)
		}
	}

	for _, beforeRequestHook := range beforeRequestHooks {
		beforeRequestHook(client, response)
	}

	if upload != nil {
		upload.upload(req, formParam)
	} else {
		if !formParam.Empty() {
			if req.Body != nil {
				//如果存在body 就特殊处理 把formData拼接到QueryString后面
				queryString.Copy(formParam)
			} else {
				setBody(req, []byte(formParam.Encode()))
				setContentType(req, "application/x-www-form-urlencoded")
			}
		}

	}

	//判断queryString不为空  拼接url
	if !queryString.Empty() {
		paramStr := queryString.Encode()
		if strings.IndexByte(reqUrl, '?') == -1 {
			reqUrl = reqUrl + "?" + paramStr
		} else {
			reqUrl = reqUrl + "&" + paramStr
		}
	}

	u, err := url.Parse(reqUrl)
	if err != nil {
		response.Err = err
		return response
	}

	req.URL = u

	if host := req.Header.Get("Host"); host != "" {
		req.Host = host
	}

	if req.Body != nil {
		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return &Response{Err: errors.New("读取Body失败")}
		}

		buf := bytes.NewBuffer(reqBody)
		req.Body = ioutil.NopCloser(buf)
		response.reqBody = reqBody
	}
	before := time.Now()
	res, err := client.client.Do(req)
	if err != nil {

		response.Err = err
		return response
	}

	response.resp = res

	after := time.Now()
	response.timeCost = after.Sub(before)

	for _, afterRequestHook := range afterRequestHooks {
		afterRequestHook(client, response)
	}

	if _, ok := response.httpClient.Transport.(*http.Transport); ok && response.resp.Header.Get("Content-Encoding") == "gzip" && req.Header.Get("Accept-Encoding") != "" {
		body, err := gzip.NewReader(res.Body)
		if err != nil {
			response.Err = err
			return response
		}
		res.Body = body
	}

	if res.Body != nil {
		respBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			response.Err = err
			return response
		}

		response.respBody = respBody
		res.Body = ioutil.NopCloser(bytes.NewBuffer(respBody))
	}

	//	logging.Info("请求结果为 --- " + response.GetBodyString())

	if retryFunc != nil {
		if retryFunc(client, response) {
			return &Response{Err: ErrCanRetry}
		}
	}

	return response
}

func setContentType(req *http.Request, contentType string) {
	req.Header.Set("Content-Type", contentType)
}

func (req *Client) Clone() *Client {
	return &Client{client: req.client}
}

func (req *Client) DoWithTimeOut(method, url string, timeout time.Duration, vals ...interface{}) *Response {
	reqCopy := req.Clone()
	reqCopy.client.Timeout = timeout

	var resChan = make(chan *Response)
	timer := time.NewTimer(timeout)

	go func() {
		resp := reqCopy.Do(method, url, vals...)
		resChan <- resp
	}()
	select {
	case resp := <-resChan:
		return resp
	case <-timer.C:
		return &Response{Err: ErrTimeOut}
	}

}

func (req *Client) DoAsync(method, url string, vals ...interface{}) <-chan *Response {
	var resChan = make(chan *Response)

	go func() {
		resp := req.Do(method, url, vals...)

		resChan <- resp
	}()

	return resChan
}

func (req *Client) DoAsyncWithTimeout(method, url string, timeout time.Duration, vals ...interface{}) <-chan *Response {
	var resChan = make(chan *Response)
	reqCopy := req.Clone()
	reqCopy.client.Timeout = timeout

	go func() {
		timer := time.NewTimer(timeout)
		var ok = make(chan struct{})
		var resp *Response

		go func() {
			resp = reqCopy.Do(method, url, vals...)
			ok <- struct{}{}
		}()

		select {
		case <-ok:
			resChan <- resp
		case <-timer.C:
			resChan <- &Response{Err: ErrTimeOut}
		}
	}()

	return resChan
}

func (req *Client) Get(url string, v ...interface{}) *Response {
	return req.Do("GET", url, v...)
}

func (req *Client) Post(url string, v ...interface{}) *Response {
	return req.Do("POST", url, v...)
}

func (req *Client) Put(url string, v ...interface{}) *Response {
	return req.Do("PUT", url, v...)
}

func (req *Client) Patch(url string, v ...interface{}) *Response {
	return req.Do("PATCH", url, v...)
}

func (req *Client) Delete(url string, v ...interface{}) *Response {
	return req.Do("DELETE", url, v...)
}

func (req *Client) Head(url string, v ...interface{}) *Response {
	return req.Do("GET", url, v...)
}

func Get(url string, v ...interface{}) *Response {
	return GetDefaultClient().Get(url, v...)
}

func Post(url string, v ...interface{}) *Response {
	return GetDefaultClient().Post(url, v...)
}

func Put(url string, v ...interface{}) *Response {
	return GetDefaultClient().Put(url, v...)
}

func Patch(url string, v ...interface{}) *Response {
	return GetDefaultClient().Patch(url, v...)
}

func Delete(url string, v ...interface{}) *Response {
	return GetDefaultClient().Delete(url, v...)
}

func Head(url string, v ...interface{}) *Response {
	return GetDefaultClient().Head(url, v...)
}

func (f *FileUpload) upload(req *http.Request, formParam param) error {

	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)
	defer writer.Close()


	for key, vals := range formParam.getValues() {
		for _, val := range vals {
			writer.WriteField(key, val)
		}
	}

	formFile, err := writer.CreateFormFile(f.FieldName, f.FileName)
	if err != nil {
		return err
	}
	_, err = io.Copy(formFile, f.File)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Body = ioutil.NopCloser(body)
	return nil
}
