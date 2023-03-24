package initialize

import (
	"github.com/issueye/lichee/app/global"
	"github.com/issueye/lichee/pkg/task"
	"github.com/spf13/cast"
)

type TaskJob struct {
	Name string
	Id   int64
	Path string
}

func (t TaskJob) Run() {
	vm := global.GetInitCore()
	err := vm.Run(t.Path)
	if err != nil {
		global.Log.Errorf("运行脚本【%s】失败，失败原因：%s", err.Error())
		return
	}
}

func InitTaskJob() {
	for _, job := range global.LocalCfg.Job {
		if job.Benable {
			task.OwnerTask.Schedule(job.Expr, cast.ToString(job.Id),
				TaskJob{
					Name: job.Name,
					Id:   job.Id,
					Path: job.Path,
				})
		}
	}
}
