package requests

import (
	"strconv"
)



var delimiter = "========================" + "\n"
func (resp *Response) Debug() {
	resp.GetDeBugInfo(resp.client.debugFlag)
}

func (resp *Response) GetDeBugInfo(flag uint8) string {

	if resp.Err != nil {
		return ""
	}
	var dumpInfo string
	//判断是否有打印头部的标志
	if flag&HeaderDump != 0 {
		headerInfo := resp.dumpHeader()
		dumpInfo += headerInfo
	}

	if flag&BodyDump != 0 {
		bodyInfo := resp.dumpBody()
		dumpInfo += bodyInfo
	}

	if flag&RespHeadDump != 0 {
		respHeadDump := resp.dumpRespHead()
		dumpInfo += respHeadDump
	}

	if flag&RespBodyDump != 0 {
		respBodyDump := resp.dumpRespBody()
		dumpInfo += respBodyDump
	}

	if flag&TimeCostDump != 0 {
		timeCostDump := resp.dumpTimeCost()
		dumpInfo += timeCostDump
	}

	return dumpInfo
}

func (resp *Response) dumpHeader() string {
	var headerInfo string

	headerInfo += resp.request.Method + " " + resp.request.URL.String() + " " + resp.request.Proto + "\n"

	for headerKey, HeaderVals := range resp.request.Header {
		for _, headerVal := range HeaderVals {
			headerInfo += headerKey + " : " + headerVal + "\n"
		}

	}

	headerInfo += delimiter
	return headerInfo

}

func (resp *Response) dumpBody() string {
	var bodyInfo string

	bodyInfo += "body: \n"
	bodyInfo += string(resp.reqBody) + " \n"
	bodyInfo += delimiter
	return bodyInfo
}

func (resp *Response) dumpRespHead() string {
	var respHeaderInfo string

	respHeaderInfo += resp.resp.Proto + " " + strconv.Itoa(resp.resp.StatusCode) + " " + resp.resp.Status + "\n"

	for headerKey, HeaderVals := range resp.resp.Header {
		for _, headerVal := range HeaderVals {
			respHeaderInfo += headerKey + " : " + headerVal + "\n"
		}

	}

	respHeaderInfo += delimiter
	return respHeaderInfo
}

func (resp *Response) dumpRespBody() string {
	return string(resp.respBody) + "\n"
}

func (resp *Response) dumpTimeCost() string {
	return "请求总耗时: " + resp.timeCost.String() + "\n"
}
