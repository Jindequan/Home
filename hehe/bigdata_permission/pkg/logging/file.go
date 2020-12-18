package logging

import (
	"fmt"
	"bigdata_permission/conf"
	"time"
)

// getLogFilePath get the log file save path
func getLogFilePath() string {
	return fmt.Sprintf("%s", conf.BaseConf.LogPath)
}

// getLogFileName get the save name of the log file
func getLogFileName() string {
	return fmt.Sprintf("%s-%s.log",
		conf.BaseConf.App,
		time.Now().Format("20060102"),
	)
}
