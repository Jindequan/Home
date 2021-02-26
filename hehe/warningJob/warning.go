package warningJob

import (
	"encoding/json"
	"time"
	"web_leads_backend/dao"
	"web_leads_backend/pkg"
	"web_leads_backend/pkg/redis"
	"web_leads_backend/service/warning"
)

func Run() {
	RefreshJob()
	DispatchJob()
}

func DispatchJob() {
	stopTime := time.Now().Unix()
	startTime := stopTime - 3600
	_, list := dao.FindCurrentWarningTask(int(startTime), int(stopTime))
	for _, task := range list{
		 go ExecTask(task) //异步，防止阻塞
	}
}

func ExecTask(task *dao.WarningTask) {
	str, _ := json.Marshal(task)
	if task.CliName == "" {
		sendToRobot(1, "未定义cli_name:" + string(str))
		return
	}
	job := JobInfo{
		JobName: task.CliName,
		JobRedisLockKey: task.CliName + "_lock",
	}
	j := GetJobInfoByName(&job)
	if j == nil {
		sendToRobot(1, "未定义Job构造方法:" + string(str))
		return
	}
	//redis lock
	locked := redis.SetNX(job.JobRedisLockKey, 3600)
	if !locked {
		return
	}
	defer redis.Del(job.JobRedisLockKey)
	//mysql record lock
	updateInfo := &dao.WarningTask{
		IsLocked: 1,
	}
	if !task.Update(updateInfo) {
		sendToRobot(1, "更新脚本数据库锁失败" + string(str))
		return
	}
	defer dao.UnLockWarningTask(task.Id)

	//run job
	startTimeStamp := time.Now().UnixNano() / 1e6
	isSuccess, errorMsg := j.Run(*task)
	stopTimeStamp := time.Now().UnixNano() / 1e6

	//save record
	record := dao.WarningRecord{
		TaskId: task.Id,
		UsedTime: uint(stopTimeStamp - startTimeStamp),
		IsSuccess: 1,
		Remark: errorMsg,
		AssignUids: task.AssignUids,
	}
	if !isSuccess {
		record.IsSuccess = 2
		sendToRobot(1, task.Name + "\n执行出错，错误信息：\n" + errorMsg)
	}
	record.Insert()

	//update task info
	_, newTask := dao.FindWarningTaskById(task.Id)
	getNewTime, nextExecTime := GetNextExecTime(newTask)
	if !getNewTime || nextExecTime == 0 {
		sendToRobot(1, "获取下次执行时间失败" + string(str))
	}
	newUpdateInfo := &dao.WarningTask{
		LastExecTime: uint(startTimeStamp / 1000),
	}
	if getNewTime {
		newUpdateInfo.NextExecTime = nextExecTime
	}
	newTask.Update(newUpdateInfo)
	return
}

func RefreshJob() {
	//初始化新脚本
	_, noTimeJob := dao.FindNoTimeTask()
	for _, task := range noTimeJob {
		if task.IsLocked == 1 {
			continue
		}
		SetNextExecTime(task)
	}
	//刷新过期时间的脚本
	_, expiredJob := dao.FindExpiredWarningTask(int(time.Now().Unix()) - 3600)
	for _, task := range expiredJob {
		if task.IsLocked == 1 {
			continue
		}
		SetNextExecTime(task)
	}
}

func SetNextExecTime(task *dao.WarningTask) (bool, uint) {
	isSuccess, nextTime := GetNextExecTime(task)
	if !isSuccess {
		return false, uint(0)
	}
	updateInfo := &dao.WarningTask{
		NextExecTime: nextTime,
	}

	return task.Update(updateInfo), nextTime
}

func GetNextExecTime(task *dao.WarningTask) (bool, uint) {
	isSuccess, nextTime := false, uint(0)

	switch task.ExecType {
	case 1:
		isSuccess, nextTime = GetNextExecTimeHourly(task)
		break
	case 2:
		isSuccess, nextTime = GetNextExecTimeDaily(task)
		break
	case 3:
		isSuccess, nextTime = GetNextExecTimeWeekly(task)
		break
	case 4:
		isSuccess, nextTime = GetNextExecTimeMonthly(task)
		break
	}
	if !isSuccess || nextTime == 0 {
		return false, uint(0)
	}
	return true, nextTime
}

func GetNextExecTimeHourly(task *dao.WarningTask) (bool, uint) {
	isOk, execTimeArray := checkExecTime(task)
	if !isOk {
		return false, uint(0)
	}

	now := time.Now()//当前时间，用来比较，处理得到的时间大于等于此值时可以更新，否则取下一个值
	nextTimeArr := []int{}//处理得到的时间必须是大于当前时间
	for _, execTime := range execTimeArray {
		min := int(execTime.Minute)
		//小于当前分钟：已过期，取下小时
		gap := 0
		if min <= now.Minute() {
			gap += 1
		}
		nowDate := now.Format("2006-01-02")
		nextTime := int(pkg.StringToTime(nowDate) / 1000) + (now.Hour() + gap) * 3600 + min * 60
		nextTimeArr = append(nextTimeArr, nextTime)
	}
	if len(nextTimeArr) == 0 {
		return false, uint(0)
	}
	return true, uint(getMinTime(nextTimeArr))
}

func GetNextExecTimeDaily(task *dao.WarningTask) (bool, uint) {
	isOk, execTimeArray := checkExecTime(task)
	if !isOk {
		return false, uint(0)
	}

	now := time.Now()//当前时间，用来比较，处理得到的时间大于等于此值时可以更新，否则取下一个值
	nextTimeArr := []int{}//处理得到的时间必须是大于当前时间
	for _, execTime := range execTimeArray {
		hour := int(execTime.Hour)
		min := int(execTime.Minute)

		gap := 0
		//小于当前小时：已过期，取下一天 || 处于当前小时但分钟过期的取下一天
		if hour < now.Hour() || (hour == now.Hour() && min <= now.Minute()) {
			gap += 1
		}

		nextDate := now.AddDate(0, 0, gap).Format("2006-01-02")
		nextTime := int(pkg.StringToTime(nextDate) / 1000) + hour * 3600 + min * 60
		nextTimeArr = append(nextTimeArr, nextTime)
	}
	if len(nextTimeArr) == 0 {
		return false, uint(0)
	}
	return true, uint(getMinTime(nextTimeArr))
}

func GetNextExecTimeWeekly(task *dao.WarningTask) (bool, uint) {
	isOk, execTimeArray := checkExecTime(task)
	if !isOk {
		return false, uint(0)
	}

	now := time.Now()//当前时间，用来比较，处理得到的时间大于等于此值时可以更新，否则取下一个值
	currentWeek := int(now.Weekday())//今天是周几
	nextTimeArr := []int{}//处理得到的时间必须是大于当前时间
	for _, execTime := range execTimeArray {
		week := int(execTime.Week)
		hour := int(execTime.Hour)
		min := int(execTime.Minute)
		if week == 7 {//以防手动插入的数据有不合规的星期
			week = 0
		}
		gap := week - currentWeek
		//如果等于0,恰好是今天，需要判断hour 与 minute
		//小于当前小时：已过期，取下周 || 处于当前小时但分钟过期的取下周
		if gap == 0 && (hour < now.Hour() || (hour == now.Hour() && min <= now.Minute())) {
			gap += 7
		}
		//如果小于0，取下周
		if gap < 0 {
			gap += 7
		}
		nextDate := now.AddDate(0, 0, gap).Format("2006-01-02")
		nextTime := int(pkg.StringToTime(nextDate) / 1000) + hour * 3600 + min * 60
		nextTimeArr = append(nextTimeArr, nextTime)
	}
	if len(nextTimeArr) == 0 {
		return false, uint(0)
	}
	return true, uint(getMinTime(nextTimeArr))
}

func GetNextExecTimeMonthly(task *dao.WarningTask) (bool, uint) {
	isOk, execTimeArray := checkExecTime(task)
	if !isOk {
		return false, uint(0)
	}

	now := time.Now()//当前时间，用来比较，处理得到的时间大于等于此值时可以更新，否则取下一个值
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	nextTimeArr := []int{}//处理得到的时间必须是大于当前时间
	for _, execTime := range execTimeArray {
		day := int(execTime.Day)
		hour := int(execTime.Hour)
		min := int(execTime.Minute)

		gapMonth := 0
		//小于当前小时：已过期，取下一天 || 处于当前小时但分钟过期的取下一天
		if day < now.Day() || (day == now.Day() && (hour < now.Hour() || (hour == now.Hour() && min <= now.Minute()))) {
			gapMonth += 1
		}

		nextMonth := int(time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation).
			AddDate(0, gapMonth,0).Unix())
		nextTime := nextMonth + (day - 1) * 24 * 3600 + hour * 3600 + min * 60
		nextTimeDay := time.Unix(int64(nextTime), 0).Day()
		if nextTimeDay != day { // 2月30日 4月31等不存在的日期
			return false, uint(0)
		}
		nextTimeArr = append(nextTimeArr, nextTime)
	}
	if len(nextTimeArr) == 0 {
		return false, uint(0)
	}
	return true, uint(getMinTime(nextTimeArr))
}

func checkExecTime(task *dao.WarningTask) (bool, []warning.ExecTimeSingle) {
	now := time.Now().Unix()
	execTimeArray := []warning.ExecTimeSingle{}
	if int(task.NextExecTime) >= int(now) {
		return true, execTimeArray
	}

	err := json.Unmarshal([]byte(task.ExecTime), &execTimeArray)
	taskJson, _ := json.Marshal(task)
	if err != nil {
		sendToRobot(1, "该任务错误的时间格式:" + string(taskJson))
		return false, execTimeArray
	}
	return true, execTimeArray
}

func sendToRobot(level int, msg string) {
	robotLevel := 0
	switch level {
	case 1:
		robotLevel = pkg.ROBOT_ERROR
		break
	case 2:
		robotLevel = pkg.ROBOT_WARNING
		break
	}

	//fmt.Println(robotLevel, "自动预警脚本任务: " + msg)
	pkg.SendToRobot(robotLevel, "自动预警脚本任务: " + msg)
}

func getMinTime(values []int) int{
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values {
		if (v < min) {
			min = v
		}
	}
	return min
}