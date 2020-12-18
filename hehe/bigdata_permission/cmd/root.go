package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)


var rootCmd = &cobra.Command{
	Short: "脚手架cli入口",
	Long:  `脚手架cli入口`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`* 　　　　　　　　┏┓　　　┏┓+ +
* 　　　　　　　┏┛┻━━━┛┻┓ + +
* 　　　　　　　┃　　　　　　　┃
* 　　　　　　　┃　　　━　　　┃ ++ + + +
* 　　　　　　 ████━████ ┃+
* 　　　　　　　┃　　　　　　　┃ +
* 　　　　　　　┃　　　┻　　　┃
* 　　　　　　　┃　　　　　　　┃ + +
* 　　　　　　　┗━┓　　　  ┏━┛
* 　　　　　　　　　┃　　　┃
* 　　　　　　　　　┃　　　┃ + + + +
* 　　　　　　　　　┃　　　┃　　　　
* 　　　　　　　　　┃　　　┃ + 　　
* 　　　　　　　　　┃　　　┃
* 　　　　　　　　　┃　　　┃　　+
* 　　　　　　　　　┃　 　　┗━━━┓ + +
* 　　　　　　　　　┃ 　　　　　　　┣┓
* 　　　　　　　　　┃ 　　　　　　　┏┛
* 　　　　　　　　　┗┓┓┏━┳┓┏┛ + + + +
* 　　　　　　　　　　┃┫┫　┃┫┫
* 　　　　　　　　　　┗┻┛　┗┻┛+ + + +`)
	},

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
