package requests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)
var client = New(nil)

func TestClient_Do_FormData(t *testing.T) {
	var param1Data = "测试测试测试"
	var param2Data = "test2222222"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "POST"{
			t.Errorf("发出的是POST请求 接受的却是 '%s'", r.Method)
		}
		//if r.URL.EscapedPath() != "/person" {
		//	t.Errorf("Expected request to '/person', got '%s'", r.URL.EscapedPath())
		//}
		r.ParseForm()
		param1 := r.Form.Get("param1")
		if param1 !=  param1Data{
			t.Errorf("发出的参数1接受错误 实际接收值为: '%s'", param1)
		}
		param2 := r.Form.Get("param2")
		if param2 !=  param2Data{
			t.Errorf("发出的参数2接受错误 实际接收值为: '%s'", param1)
		}
	}))

	client.Do("POST",server.URL,FormData{"param1":param1Data,"param2":param2Data})
}

func TestClient_DoQueryString(t *testing.T) {
	var param1Data = "测试测试测试"
	var param2Data = "test2222222"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "GET"{
			t.Errorf("发出的是GET请求 接受的却是 '%s'", r.Method)
		}
		//if r.URL.EscapedPath() != "/person" {
		//	t.Errorf("Expected request to '/person', got '%s'", r.URL.EscapedPath())
		//}
		r.ParseForm()
		param1 := r.Form.Get("param1")
		if param1 !=  param1Data{
			t.Errorf("发出的参数1接受错误 实际接收值为: '%s'", param1)
		}
		param2 := r.Form.Get("param2")
		if param2 !=  param2Data{
			t.Errorf("发出的参数2接受错误 实际接收值为: '%s'", param1)
		}
	}))

	client.Do("GET",server.URL,QueryString{"param1":param1Data,"param2":param2Data})
}

func TestClient_DoString(t *testing.T) {
	var str = "对双方都是打工大锅饭大锅饭符合规范放过机会了发卢浮宫黄金分割厉害就够了好几个来回就更好了机构行家律师代理费收到反馈速度的付款了国际饭店旅馆附件改好了房间和法国黄金分割"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error("读取出错")
		}

		if str != string(body) {
			t.Errorf("实际接收值为:%s错误",string(body))
		}

		if r.Method != "POST"{
			t.Errorf("发出的是POST请求 接受的却是 '%s'", r.Method)
		}
	}))

	client.Do("POST",server.URL,str)
}
func TestClient_DoJsonBody(t *testing.T) {
	type TestJson struct {
		W int `json:"w"`
		Q string `json:"q"`
	}

	var testJson = &TestJson{1,"test"}
	jsonBody,_ := json.Marshal(testJson)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		body,err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error("读取出错")
		}


		var recive = &TestJson{}
		json.Unmarshal(body,recive)
		if !reflect.DeepEqual(recive,testJson) {
			t.Error("接收数据不匹配")
		}
		w.Write(body)

		if r.Method != "POST"{
			t.Errorf("发出的是POST请求 接受的却是 '%s'", r.Method)
		}
	}))

	resp := client.Do("POST",server.URL,JsonBodyStr(jsonBody))
	decodeRes := &TestJson{}
	resp.GetJson(decodeRes)

	if !reflect.DeepEqual(decodeRes,testJson) {
		t.Error("GetJson方法不匹配")
	}
}

func TestNew(t *testing.T) {
	client := new(Client)

	newDefaultClient := New(nil)

	newClient := New(&http.Client{})

	if reflect.TypeOf(client) != reflect.TypeOf(newDefaultClient) {
		t.Fatal("New不传参数方法创建对象失败")
	}

	if reflect.TypeOf(client) != reflect.TypeOf(newClient) {
		t.Fatal("New传递参数方法创建对象失败")
	}

	t.Log("New方法sucess")

}

func TestClient_DoWithTimeOut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Microsecond)

	}))

	resp := client.DoWithTimeOut("GET",server.URL,80 * time.Microsecond)
	if resp.Err != ErrTimeOut {
		t.Error("超时机制失败")
	}

	resp = client.DoWithTimeOut("GET",server.URL,1 * time.Second)
	if resp.Err == ErrTimeOut {
		t.Error("？？？")
	}

}
