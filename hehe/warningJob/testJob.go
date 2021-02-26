package warningJob

import (
	"fmt"
	"web_leads_backend/dao"
)

func (job *TestJob) Run(task dao.WarningTask) (bool, string) {
	fmt.Println("here is test job executable function")
	return false, ""
}