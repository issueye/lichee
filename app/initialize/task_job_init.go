package initialize

import (
	"time"

	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/app/model"
	"github.com/issueye/lichee/app/service"
	"github.com/issueye/lichee/global"
	"github.com/issueye/lichee/pkg/plugins/core"
	"github.com/issueye/lichee/utils"
	"github.com/spf13/cast"
	"go.etcd.io/bbolt"
)

type TaskJob struct {
	Name   string
	Id     int64
	Path   string
	AreaId int64
}

func (t TaskJob) Run() {
	vm := common.GetInitCore()

	// 注入系统参数
	err := regParam(common.SYS_AREA, vm)
	if err != nil {
		common.Log.Errorf("系统参数注入失败，失败原因：%s", err.Error())
		return
	}

	// 注入脚本对应的参数域
	err = regParam(t.AreaId, vm)
	if err != nil {
		common.Log.Errorf("参数注入失败，失败原因：%s", err.Error())
		return
	}

	err = vm.Run(t.Name, t.Path)
	if err != nil {
		common.Log.Errorf("运行脚本【%s】失败，失败原因：%s", t.Path, err.Error())
		return
	}
}

func regParam(id int64, vm *core.Core) error {
	pa, err := service.NewParamService().GetAreaById(id)
	if err != nil {
		common.Log.Errorf("根据编码【%d】未找到参数域信息", id)
		return err
	}

	// 注入参数域的参数
	global.Bdb.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(common.AREA_BUCKET)
		paramBucket := b.Bucket(common.AreaBucketID(pa.Id))
		err := paramBucket.ForEach(func(k, v []byte) error {
			if len(v) > 0 {
				data := new(model.Param)
				err := utils.GobBuff{}.BytesToStruct(v, data)
				if err != nil {
					return err
				}

				vm.SetProperty(pa.Name, data.Name, data.Value)
			}

			return nil
		})
		return err
	})

	return nil
}

func InitTaskJob() {
	go Monitor()
	time.Sleep(200 * time.Millisecond)
	list, err := service.NewJobService().Query(new(model.ReqQueryJob))
	if err != nil {
		return
	}

	for _, job := range list {
		if job.Enable && job.AreaId != 0 {
			Job := model.Job{
				Id:     job.Id,
				Name:   job.Name,
				Expr:   job.Expr,
				Mark:   job.Mark,
				Enable: job.Enable,
				AreaId: job.AreaId,
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
							Name:   job.Name,
							Id:     job.Id,
							Path:   job.Path,
							AreaId: job.AreaId,
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
						// 移除日志对象
						zl, ok := core.LogMap[job.Name]
						if ok {
							zl.Close()
							delete(core.LogMap, job.Name)
						}
					}
				case common.JOB_MODIFY: // 修改定时任务  先移除任务 再添加定时任务
					{
						global.JobTask.Remove(cast.ToString(job.Id))

						// 添加定时任务
						tj := TaskJob{
							Name:   job.Name,
							Id:     job.Id,
							Path:   job.Path,
							AreaId: job.AreaId,
						}
						_, err := global.JobTask.Schedule(job.Expr, cast.ToString(job.Id), tj)
						if err != nil {
							common.Log.Errorf("添加定时任务失败，失败原因：%s", err.Error())
							return
						}

						common.Log.Debugf("定时任务【%s-%d】修改成功", tj.Name, tj.Id)
					}
				case common.JOB_AT_ONCE_RUN:
					{
						tj := TaskJob{
							Name:   job.Name,
							Id:     job.Id,
							Path:   job.Path,
							AreaId: job.AreaId,
						}
						// 马上运行定时任务
						global.JobTask.Now(tj)
					}
				case common.JOB_DELAY_ONCE_RUN:
					{
						tj := TaskJob{
							Name:   job.Name,
							Id:     job.Id,
							Path:   job.Path,
							AreaId: job.AreaId,
						}
						// 延迟时间运行一次
						global.JobTask.In(job.Delay, tj)
					}
				}

			}
		}
	}
}
