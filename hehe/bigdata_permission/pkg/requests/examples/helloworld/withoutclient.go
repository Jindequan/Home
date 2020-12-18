package main

import (
	"fmt"
	"account_assets/pkg/requests"
)

type response struct {
	code int
	msg string
}
func main() {

	client := requests.New(nil)

	resp := client.Do("GET","http://www.baidu.com/s",requests.QueryString{"wd":"test2"},requests.FormData{"testForm":1})
	fmt.Println(resp.GetBodyString())
	fmt.Println(resp.GetBodyBytes())
	var res = response{}
	fmt.Println(resp.GetJson(&res))
}
