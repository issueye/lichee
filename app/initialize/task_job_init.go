package initialize

import (
	"time"

	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/app/model"
	"github.com/issueye/lichee/app/service"
	"github.com/issueye/lichee/global"
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
	go Monitor()
	time.Sleep(200 * time.Millisecond)
	list, err := service.NewJobService().Query(new(model.ReqQueryJob))
	if err != nil {
		return
	}

	for _, job := range list {
		if job.Enable {
			Job := model.Job{
				Id:     job.Id,
				Name:   job.Name,
				Expr:   job.Expr,
				Mark:   job.Mark,
				Enable: job.Enable,
				Path:   job.Path,
			}

			common.JobGo(Job, common.JOB_ADD)
		}
	}
}

func Monitor() {
	for {
		select {
		case job := <-common.TASK_CHAN:
			{
				switch job.Type {
				case common.JOB_ADD: // 添加定时任务
					{
						tj := TaskJob{
							Name: job.Name,
							Id:   job.Id,
							Path: job.Path,
						}
						_, err := global.JobTask.Schedule(job.Expr, cast.ToString(job.Id), tj)
						if err != nil {
							common.Log.Errorf("添加定时任务失败，失败原因：%s", err.Error())
							return
						}

						common.Log.Debugf("定时任务【%s-%d】添加成功", tj.Name, tj.Id)
					}
				case common.JOB_DEL: // 删除定时任务
					{
						global.JobTask.Remove(cast.ToString(job.Id))
						common.Log.Debugf("定时任务【%s-%d】删除成功", job.Name, job.Id)
					}
				case common.JOB_MODIFY: // 修改定时任务  先移除任务 再添加定时任务
					{
						global.JobTask.Remove(cast.ToString(job.Id))

						// 添加定时任务
						tj := TaskJob{
							Name: job.Name,
							Id:   job.Id,
							Path: job.Path,
						}
						_, err := global.JobTask.Schedule(job.Expr, cast.ToString(job.Id), tj)
						if err != nil {
							common.Log.Errorf("添加定时任务失败，失败原因：%s", err.Error())
							return
						}

						common.Log.Debugf("定时任务【%s-%d】修改成功", tj.Name, tj.Id)
					}
				}
			}
		}
	}
}
