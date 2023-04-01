package initialize

import (
	"fmt"

	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/pkg/task"
	"github.com/issueye/lichee/utils"
	"github.com/spf13/cast"
)

type TaskJob struct {
	Name string
	Id   int64
	Path string
}

func (t TaskJob) Run() {
	vm := common.GetInitCore()
	err := vm.Run(t.Path)
	if err != nil {
		common.Log.Errorf("运行脚本【%s】失败，失败原因：%s", err.Error())
		return
	}
}

func InitTaskJob() {
	for _, job := range common.LocalCfg.Job {
		if job.Benable {
			task.OwnerTask.Schedule(job.Expr, cast.ToString(job.Id),
				TaskJob{
					Name: job.Name,
					Id:   job.Id,
					Path: job.Path,
				})
		}
	}
	fmt.Printf("【%s】初始化定时任务完成...\n", utils.Ltime{}.GetNowStr())
}
