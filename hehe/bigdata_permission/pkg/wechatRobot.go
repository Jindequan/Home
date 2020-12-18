package pkg

import (
	"bigdata_permission/pkg/requests"
	"encoding/json"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"time"
)

const (
	//url
	ROBOT_URL     = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?debug=1&key="
	//error level
	ROBOT_ERROR   = 1
	ROBOT_WARNING = 2
	ROBOT_INFO = 4
)

var keyMap = map[string]string {
	"online_error":   "31504a3d-b3f1-428a-b05f-dbb492a02851",
	"online_info": 	"4581f53f-b6d1-4c19-b953-914c341eef01",
	"test_all": "e08f77ba-2494-4405-a19a-da2be2113a85",
}

type textRobotBody struct {
	Msgtype string `json:"msgtype"`
	Text text `json:"text"`
}

type text struct{
	Content string `json:"content"`
}
func SendToRobot(errorLevel int, message string) bool {

	key := getRobotId(errorLevel)
	url := ROBOT_URL + key

	msg := genRobotMsg(errorLevel)

	data := textRobotBody {
		Msgtype: "text",
		Text:    text{msg + message},
	}
	jsonBody, _ := json.Marshal(data)
	response := requests.Post(url, requests.JsonBodyStr(jsonBody))

	responseBody := response.GetBodyString()
	if gjson.Get(responseBody, "errcode").Int() != 0 {
		return false
	}
	return true
}


func getRobotId(errorLevel int) string {
	if IsProduction() {
		return getProductionRobot(errorLevel)
	}
	return keyMap["test_all"]
}

func getProductionRobot(errorLevel int) string {
	switch errorLevel {
	case ROBOT_ERROR:
		return keyMap["online_error"]
	}
	return keyMap["online_info"]
}

func genRobotMsg(errorLevel int) string {
	env := strings.ToUpper(GetEnv())
	timeStr := TimeStampToString(time.Now().Unix())
	msg := "BigDataPermission"
	switch errorLevel {
	case ROBOT_ERROR:
		msg += "【ERROR】"
		break
	case ROBOT_WARNING:
		msg += "【WARNING】"
		break
	case ROBOT_INFO:
		msg += "【INFO】"
		break
	default:
		msg += "【" + strconv.Itoa(errorLevel) + "】"
		break
	}
	msg += "【" + env + "】【" + timeStr + "】\n"
	return msg
}

