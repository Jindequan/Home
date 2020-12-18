package main

import (
	"fmt"
	"account_assets/pkg/requests"
)

type test struct {
	a int
}
func main()  {
	//直接简单调用 不需要创建client对象
	resp := requests.Get("http://www.baidu.com",requests.QueryString{"a":1})
	//resp := requests.Post("www.baidu.com",requests.FormData{"b":2})
	//resp := requests.Patch("www.baidu.com",requests.JsonBodyStr("{\"a\":1}"))
	//resp := requests.Head("www.baidu.com")
	//resp := requests.Delete("www.baidu.com")
	if resp.Err != nil {
		panic(resp.Err)
	}

	fmt.Println(resp.GetBodyString())
	//fmt.Println(resp.GetBodyBytes())
	//var t test
	//fmt.Println(resp.GetJson(&t))

}