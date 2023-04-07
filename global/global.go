package global

import (
	gdb "github.com/issueye/lichee/pkg/db"
	"github.com/issueye/lichee/pkg/plugins/core/boltdb"
	"github.com/issueye/lichee/pkg/task"
	"gorm.io/gorm"
)

// 将全局变量都引入到此单元中
var (
	JobTask = task.OwnerTask
	Bdb     = boltdb.Bdb
)

type DbInfo struct {
	Name string
	Cfg  *gdb.Config
	DB   *gorm.DB
}

var (
	GdbMap = make(map[string]*DbInfo)
)
