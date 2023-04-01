//go:generate goversioninfo
package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/issueye/lichee/app/common"
	"github.com/issueye/lichee/app/initialize"
	_ "github.com/issueye/lichee/docs"
)

var (
	config = flag.String("conf", "config.json", "配置文件")
)

func main() {
	flag.Parse()

	common.ConfigPath = *config
	initialize.Initialize()
	wait()
}

func wait() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
