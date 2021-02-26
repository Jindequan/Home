package warningJob

import (
	"web_leads_backend/dao"
)

type JobInfo struct {
	JobName string
	JobRedisLockKey string
}

type JobRun interface {
	Run(task dao.WarningTask) (bool, string) //是否成功, 失败消息
}
//=============注册脚本结构体==================//
type TestJob struct {
	*JobInfo
}
type LixiaoUpdateJob struct {
	*JobInfo
}
//=============注册脚本结构体构造方法=====================//
func GetJobInfoByName(data *JobInfo) JobRun {
	switch data.JobName {
	case "test":
		return &TestJob{data}
	case "lixiao_update":
		return &LixiaoUpdateJob{data}
	}
	return nil
}
//================实现具体的脚本方法====================//
//一个脚本一个文件，避免脚本过多导致文件过大
