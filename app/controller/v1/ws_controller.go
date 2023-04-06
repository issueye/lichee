package v1

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/pkg/res"
	"github.com/issueye/lichee/pkg/ws"
	"github.com/issueye/tail"
)

// var monitorMap = make(map[string][]*MonitorClient)

type MonitorClient struct {
	WsClientId string
	Name       string
}

func WsLogMonitor(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		common.Log.Errorf("name不能为空")
		res.FailByMsg(ctx, "name不能为空")
		return
	}

	groupName := fmt.Sprintf("%s:%s", common.WS_LOG_GROUP, name)
	// 添加到日志组中
	id := ws.WebsocketManager.WsClient(ctx, groupName)

	// 判断是否已经打开监听
	go Monitor(groupName, id, fmt.Sprintf("%s/%s.log", "runtime/logs", name))
}

func Monitor(groupName, id string, fileName string) {
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	}
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
		return
	}

	var (
		line *tail.Line
		ok   bool
	)

	for {
		select {
		case line, ok = <-tails.Lines:
			{
				//遍历chan，读取日志内容
				if !ok {
					fmt.Printf("tail file close reopen, filename:%s\n", tails.Filename)
					time.Sleep(time.Second)
					continue
				}

				// 发送到ws 客户端
				ws.WebsocketManager.Send(id, groupName, []byte(line.Text))
			}
		case <-ws.UnRegChan:
			{
				goto title
			}
		}
	}

title:
	common.Log.Debugf("【%s】关闭日志监听", groupName)
}
