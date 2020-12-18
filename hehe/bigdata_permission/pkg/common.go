package pkg

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var env string

var TimeLayoutWithZone = "2006-01-02T15:04:05.000+08:00"
var TimeLayoutFull = "2006-01-02 15:04:05"
var TimeLayoutYmdHi = "2006-01-02 15:04"
var TimeLayoutYmd = "2006-01-02"

type Enviroment struct{
	ENV string `json:"ENV"`
}
var location, _ = time.LoadLocation("Asia/Shanghai")

func NowUnixMs() int64 {
	return time.Now().In(location).UnixNano() / 1e6
}

func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

//反射 性能低 慎用～
func InArray(value interface{}, arr interface{}) bool {
	switch reflect.TypeOf(arr).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(arr)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(value, s.Index(i).Interface()) {
				return true
			}
		}

	}

	return false
}

func StrFirstToUpper(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArry := []rune(str)
	if strArry[0] >= 97 && strArry[0] <= 122  {
		strArry[0] -=  32
	}
	return string(strArry)
}

func GetEnv() string {
	if env == "" {
		load()
	}

	return strings.ToLower(env)
}

func IsDev() bool {
	return GetEnv() == "dev"
}

func IsTest() bool {
	return GetEnv() == "test"
}

func IsProduction() bool {
	return GetEnv() == "online"
}

func DoNothing()  {
	return
}

func load()  {
	path := "./settings/config.json"

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var envStruct Enviroment

	err = json.Unmarshal(bs, &envStruct)
	if err != nil {
		panic(err)
	}

	env = envStruct.ENV
}

func StringToTime(timeString string) int64 {
	timeStampInt := int64(0)
	timeLen := len(timeString)
	var timeStampTime time.Time
	var err error

	if timeLen == 29 {
		timeStampTime, err = time.ParseInLocation(TimeLayoutWithZone, timeString, time.Local)
	} else if timeLen == 16 {
		timeStampTime, err = time.ParseInLocation(TimeLayoutYmdHi, timeString, time.Local)
	} else if timeLen == 10 {
		timeStampTime, err = time.ParseInLocation(TimeLayoutYmd, timeString, time.Local)
	} else {
		timeStampTime, err = time.ParseInLocation(TimeLayoutFull, timeString, time.Local)
	}

	if err == nil {
		timeStampInt = timeStampTime.Unix() * 1000
	}

	return timeStampInt
}

func TimeToString(timestampInt int64) string {
	if timestampInt == 0 {
		return ""
	}
	return time.Unix(timestampInt/1000, 0).In(location).Format(TimeLayoutFull)
}

func TimeStampToString(timestampInt int64) string {
	return time.Unix(timestampInt, 0).In(location).Format(TimeLayoutFull)
}

func TimeToStringWithZone(timestampInt int64) string {
	if timestampInt == 0 {
		return ""
	}
	return time.Unix(timestampInt/1000, 0).In(location).Format(TimeLayoutWithZone)
}

func DateNow() string {
	return time.Now().Format("20060102")
}

func TimeToDateString(timestamp int64) string {
	return time.Unix(timestamp/1000, 0).In(location).Format(TimeLayoutYmd)
}

func Md5(str string) string {
	data := []byte(str)
	sum := md5.Sum(data)
	return fmt.Sprintf("%x", sum)
}

func JoinIntArrToString(intArr []int, sep string) string {
	if len(intArr) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for _, intItem := range intArr {
		buffer.WriteString(sep + strconv.Itoa(intItem))
	}
	buffer.WriteString(sep)
	return strings.Trim(buffer.String(), sep)
}

func JoinStrArrToString(stringArr []string, sep string) string {
	if len(stringArr) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for _, str := range stringArr {
		buffer.WriteString(sep + str)
	}
	buffer.WriteString(sep)
	return buffer.String()
}


func ToMultiIntArr(indexArrStr string, sep string) []int {
	if len(indexArrStr) == 0 {
		return []int{}
	}

	indexSplits := strings.Split(indexArrStr, sep)

	var indexArr []int
	for _, indexStr := range indexSplits {
		index, err := strconv.Atoi(indexStr)
		if err != nil || index == 0 {
			continue
		}
		indexArr = append(indexArr, index)
	}
	return indexArr
}