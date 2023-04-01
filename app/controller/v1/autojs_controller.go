package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dop251/goja"
	"github.com/gin-gonic/gin"
	"github.com/issueye/lichee/app/common"
	mhttp "github.com/issueye/lichee/pkg/plugins/core/net/http"
)

type AutoJsController struct {
	ScriptTimeoutSec int
}

func (gszyy *AutoJsController) AutoJsReceiveServer(ctx *gin.Context) {
	rtCore := common.GetInitCore()
	rt := rtCore.GetRts()
	rt.Set("response", mhttp.NewResponse(rt, ctx.Writer))
	rt.Set("request", mhttp.NewRequest(rt, ctx.Request))

	normalEndCh := make(chan bool)
	go func() {
		err := rtCore.RunVM("index.js", rt)
		if err != nil {
			switch err := err.(type) {
			case *goja.Exception:
				ctx.Writer.WriteHeader(http.StatusInternalServerError)
				fmt.Println("*goja.Exception:", err.String())
				if v := err.Value().ToObject(rt).Get("nativeType"); v != nil {
					fmt.Printf("%T:%[1]v\n", v.Export())
				}
			case *goja.InterruptedError:
				fmt.Println("*goja.InterruptedError:", err.String())
			default:
				fmt.Println("default err:", err)
			}
			normalEndCh <- false
			return
		}

		// if This.writeResultValue {
		// 	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 	w.Write([]byte((ret).String()))
		// }
		normalEndCh <- true
	}()

	if gszyy.ScriptTimeoutSec > 0 {
		timeoutCh := time.After(time.Duration(gszyy.ScriptTimeoutSec) * time.Second)
		select {
		case <-timeoutCh:
			rt.Interrupt("run code timeout, halt")
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
		case <-normalEndCh:
		}
	} else {
		status := <-normalEndCh
		// 如果失败则直接返回失败
		if !status {
			ctx.JSON(http.StatusOK,
				map[string]any{
					"success": "false",
					"message": "失败",
				})
		}
	}
}
