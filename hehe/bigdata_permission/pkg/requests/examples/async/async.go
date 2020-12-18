package main

import (
	"fmt"
	"account_assets/pkg/requests"
)

func main()  {
	//异步请求
	client := requests.New(nil)
	resChan := client.DoAsync("GET","http://www.baidu.com")

	fmt.Println("啦啦啦 我比响应早")
	resp := <-resChan

	if resp.Err != nil{
		panic(resp.Err)
	}

	fmt.Println(resp.GetBodyString())
}
