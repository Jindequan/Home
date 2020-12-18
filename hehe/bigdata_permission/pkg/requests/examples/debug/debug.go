package main

import (
	"fmt"
	"account_assets/pkg/requests"
)

func main()  {
	//获取debug信息
	client := requests.New(nil)
	resp := client.Do("GET","http://www.baidu.com",requests.QueryString{"wd":"test1","test2":13},requests.FormData{"testForm":1})

	if resp.Err != nil{
		panic(resp.Err)
	}

	requests.Debug = true
	//可以根据需要去掉flag中的标志
	debugInfo := resp.GetDeBugInfo(requests.HeaderDump|requests.BodyDump|requests.RespHeadDump|requests.RespBodyDump|requests.TimeCostDump) //可以根据这个要打印的debug信息

	fmt.Println(debugInfo)
}
