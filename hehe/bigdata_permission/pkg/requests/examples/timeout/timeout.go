package main

import (
	"fmt"
	"account_assets/pkg/requests"
	"time"
)

func main() {
	c := requests.New(nil)
	resp := c.DoWithTimeOut("GET", "http://www.baidu.com", 2*time.Microsecond)
	//resp := c.DoAsyncWithTimeout("GET","http://www.baidu.com",5*time.Second)
	if resp.Err == requests.ErrTimeOut {
		fmt.Println(resp.Err)
	} else if resp.Err != nil {
		panic(resp.Err)
	}

	fmt.Println(resp.GetBodyString())

	resChan := c.DoAsyncWithTimeout("GET","http://www.baidu.com",2*time.Second)
	fmt.Println("========")
	resp = <-resChan
	if resp.Err == requests.ErrTimeOut {
		fmt.Println(resp.Err)
	} else if resp.Err != nil {
		panic(resp.Err)
	}

	fmt.Println(resp.GetBodyString())

}
