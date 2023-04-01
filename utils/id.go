package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
)

type Lid struct{}

var node *snowflake.Node
var lock = new(sync.Mutex)

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	// 格式化 1月2号下午3时4分5秒  2006年
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		fmt.Println(err)
		return
	}

	snowflake.Epoch = st.UnixNano() / 1e6
	node, err = snowflake.NewNode(machineID)
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

// GenID
// 生成 64 位的 雪花 ID
func (l Lid) GenID() int64 {
	lock.Lock()
	defer lock.Unlock()
	return node.Generate().Int64()
}

func init() {
	if err := Init(carbon.Now().ToDateString(), 1); err != nil {
		fmt.Println("Init() failed, err = ", err)
		return
	}
}

func (l Lid) GetUUID() string {
	return uuid.New().String()
}

func (l Lid) GetCryptId() string {
	return Sha1(l.GetUUID())
}
