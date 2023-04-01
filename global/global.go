package global

import (
	"github.com/issueye/lichee/pkg/plugins/core/boltdb"
	"github.com/issueye/lichee/pkg/plugins/core/db"
	"github.com/issueye/lichee/pkg/task"
)

// 将全局变量都引入到此单元中
var (
	JobTask = task.OwnerTask
	Bdb     = boltdb.Bdb
	Ldb     = db.Ldb
)
