package warningJob

import (
	"os"
	"runtime"
	"web_leads_backend/conf"
)

func getFilePath(name string) string {
	if runtime.GOOS == "windows" {
		return "D:\\" + name
	}
	return conf.BaseConf.TempFilePath + name
}

func checkAndCleanFile(name string) bool {
	if checkFileExist(name) && !cleanFile(name) {
		return false
	}
	return true
}

func checkFileExist(name string) bool {
	_, err := os.Stat(getFilePath(name))
	if err == nil {
		return true
	}
	if os.IsExist(err) {
		return true
	}
	return false
}

func cleanFile(name string) bool {
	return os.Remove(getFilePath(name)) == nil
}
