package task

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type Test struct {
	Name string
}

func (t Test) Run() {
	fmt.Println(t.Name, "任务执行", time.Now().Format("2006-01-02 15:04:05.000"))
}

func TestStatusJson(t *testing.T) {
	var err error
	_, err = OwnerTask.Schedule("0/2 * * * * ?", "123", Test{Name: "测试任务-StatusJson"})
	if err != nil {
		t.Error(err)
	}

	time.Sleep(5 * time.Second)
	data := OwnerTask.StatusJson()

	var jsonData []byte
	jsonData, err = json.Marshal(data)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("jsonData", string(jsonData))
}

func TestFindTask(t *testing.T) {
	var err error
	_, err = OwnerTask.Schedule("0/2 * * * * ?", "123", Test{Name: "测试任务-FindTask"})
	if err != nil {
		t.Error(err)
	}

	time.Sleep(5 * time.Second)
	ok := OwnerTask.FindTask("123")
	if ok {
		fmt.Println("定时任务存在")
	} else {
		fmt.Println("定时任务不存在")
	}
}

func TestStop(t *testing.T) {
	var err error
	_, err = OwnerTask.Schedule("0/2 * * * * ?", "123", Test{Name: "测试任务-Stop"})
	if err != nil {
		t.Error(err)
	}

	time.Sleep(5 * time.Second)
	fmt.Println("停用定时任务功能", time.Now().Format("2006-01-02 15:04:05.000"))
	OwnerTask.Stop()
	time.Sleep(5 * time.Second)
	fmt.Println("启用定时任务功能", time.Now().Format("2006-01-02 15:04:05.000"))
	OwnerTask.Start()
	time.Sleep(5 * time.Second)
}

func TestIn(t *testing.T) {
	OwnerTask.In(3*time.Second, Test{Name: "测试任务-In"})
	time.Sleep(10 * time.Second)
}

func TestRemove(t *testing.T) {
	var err error
	_, err = OwnerTask.Schedule("0/2 * * * * ?", "123", Test{Name: "测试任务-Remove"})
	if err != nil {
		t.Error(err)
	}

	time.Sleep(5 * time.Second)
	fmt.Println("移除任务", time.Now().Format("2006-01-02 15:04:05.000"))
	// 移除任务
	OwnerTask.Remove("123")

	time.Sleep(3 * time.Second)
	fmt.Println("结束", time.Now().Format("2006-01-02 15:04:05.000"))
}

func TestRunNowTask(t *testing.T) {
	OwnerTask.Now(Test{Name: "测试任务马上执行"})
	time.Sleep(1 * time.Second)
}

func TestEvery(t *testing.T) {
	OwnerTask.Every(3*time.Second, "123", Test{Name: "测试任务-Every"})
	time.Sleep(10 * time.Second)
}

func TestSchedule(t *testing.T) {
	var err error
	_, err = OwnerTask.Schedule("0/2 * * * * ?", "123", Test{Name: "测试任务-Schedule"})
	if err != nil {
		t.Error(err)
	}
	_, err = OwnerTask.Schedule("0/2 * * * * ?", "123", Test{Name: "测试任务2-Schedule"})
	if err != nil {
		t.Error(err)
	}
	time.Sleep(10 * time.Second)
}
