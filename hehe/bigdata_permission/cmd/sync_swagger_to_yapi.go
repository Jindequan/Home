package cmd

import (
	"github.com/spf13/cobra"
	"bigdata_permission/pkg/requests"
	"bigdata_permission/conf"
	"bigdata_permission/pkg/file"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const yapiUrl = "https://web-api.***.net/"

const docPath = "./docs/"
const swaagerFilePath = docPath + "swagger.json" //todo 改成配置项
const bakPath = docPath + "bak/"                 //todo 改成配置项

func init() {
	rootCmd.AddCommand(syncSwaggerApiCommand)
}

var syncSwaggerApiCommand = &cobra.Command{
	Use:   "syncdoc",//命令行名称 cli syncdoc
	Short: "生成swagger文档并同步到yapi",  //命令短注释
	Long:  "生成swagger文档并同步到yapi",  //命令长注释

	Run: func(cmd *cobra.Command, args []string) {//命令执行方法体
		t := time.Now()
		todayStr := t.Format("20060102_150405")
		date := strings.Split(todayStr, "_")[0]

		backUpFileName := bakPath + "." + todayStr + "_swagger.json"

		err := file.CopyFile(swaagerFilePath, backUpFileName)
		if err != nil {
			cmd.Printf("备份失败 %v", err)
			return
		}

		fds, err := ioutil.ReadDir(bakPath)

		if err != nil {
			cmd.Printf("清除过期swaggerjson失败 %v", err)
			return
		}
		reg, err := regexp.Compile(`\.(\d{8})_(\d{6})_swagger.json`)
		if err != nil {
			cmd.Printf("编译正则失败 %v", err)
			return
		}

		for _, fd := range fds {
			fileName := fd.Name()


			res := reg.FindSubmatch([]byte(fileName))
			if len(res) > 0 {
				createDateToInt, err := strconv.Atoi(string(res[1]))

				if err != nil {
					cmd.Printf("createDate转换失败 %v", err)
					return
				}

				todayDateToInt, err := strconv.Atoi(date)
				if err != nil {
					cmd.Printf("todayDate转换失败 %v", err)
				}

				if createDateToInt < todayDateToInt {
					err := os.Remove(bakPath + fileName)
					if err != nil {
						cmd.Println(err)
					}
				}
			}
		}

		command := exec.Command("swag", "init")
		err = command.Run()
		if err != nil {
			cmd.Printf("执行swag init失败 %v",err)
			return
		}

		inputSwaggerFile, err := ioutil.ReadFile(swaagerFilePath) //todo 做成命令行配置项

		if err != nil {
			cmd.Printf("读取swaager文件出错 %v", err)
			return
		}

		resp := requests.Post(yapiUrl+"api/open/import_data", requests.FormData{"type": "swagger", "token": conf.BaseConf.SwaagerToken, "merge": "merge", "json": string(inputSwaggerFile)})
		if resp.Err != nil {
			cmd.Printf("同步到yapi出错 %v", resp.Err)
			return
		}

		cmd.Println(resp.GetBodyString())
	},
}
