package initialize

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/app/model"
	"github.com/issueye/lichee/utils"
)

func InitConfig() {
	config := common.ConfigPath
	// 读取配置文件地址
	path, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("获取当前程序路径失败，失败原因：%s", err.Error()))
	}

	sysType := runtime.GOOS
	var data []byte

	if sysType == "linux" {
		data, err = os.ReadFile(fmt.Sprintf("%s/%s", path, config))
		if err != nil {
			panic(fmt.Errorf("获取配置文件信息失败，失败原因：%s", err.Error()))
		}
	}

	if sysType == "windows" {
		data, err = os.ReadFile(fmt.Sprintf("%s\\%s", path, config))
		if err != nil {
			panic(fmt.Errorf("获取配置文件信息失败，失败原因：%s", err.Error()))
		}
	}

	common.LocalCfg = new(model.Config)
	err = json.Unmarshal(data, common.LocalCfg)
	if err != nil {
		panic(fmt.Errorf("解析配置文件失败，失败原因：%s", err.Error()))
	}

	fmt.Printf("【%s】配置文件加载完成...\n", utils.Ltime{}.GetNowStr())
}
