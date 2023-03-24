#  Task 定时任务模块

![图标](./images/任务.png)

> 在业务中存在需要按照固定时间循环执行的任务，通过时间表达式来进行调度和执行的一类任务被称为定时任务，很多业务需求的实现都离不开定时任务，例如备忘录提醒、闹钟等功能

## 时间表达式

> 时间表达式：用来控制任务的调度与执行，时间表达式又称为：cron表达式

    cron表达式用于配置cronTrigger的实例。cron表达式实际上是由六个子表达式组成。这些表达式之间用空格分隔，如 ：" * * * * ? * "

从前到后依次代表：
1. Seconds （秒）
2. Minutes（分）
3. Hours（小时）
4. Day-of-Month （天）
5. Month（月）
6. Day-of-Week （周）

例如 ：“ * 30 12 5 * * ? ” 就代表每个月的5号中午12点30执行

Cron表达式的格式：

    秒 分 时 日 月 周 年(可选)

    字段名	允许的值	允许的特殊字符
    秒	0-59	, - * /
    分	0-59	, - * /
    小时	0-23	, - * /
    日	1-31	, - * ? / L W C
    月	1-12 or JAN-DEC	, - * /
    周	1-7 or SUN-SAT	, - * ? / L C #
    字符含义：

    " * "：代表所有可能的值。在哪个字段中代表哪个含义

    " - "：表示指定范围。

    " , " ：表示可有多个值。例如：“ * 10,20 12 * * * ? * ”，“10,20”代表在12点10分和12点20分都会触发。

    " / " ：表示执行过程。例如：“ * 10/20 12 * * * ? * ”，表示在每天的12点钟这个小时内从10分钟开始，每20分钟执行一次

    " L" ：在月字段中，" L" 表示一个月的最后一天；在周字段中，" L"表示一个星期的最后一天

    " ? " ：与 " * "含义类似，但为了避免冲突，一般在日字段 和周字段 使用时 需要将另外一个的值设为 " ? "


## 模块说明

- ### 全局对象

模块中对外提供了一个全局对象 `OwnerTask`，以便进行定时任务业务操作和管理

```golang
// OwnerTask 全局对象
var OwnerTask ICronTaskMana

type ICronTaskMana interface {
	// Start 初始化功能
	Start(v ...int)
	// Schedule 添加任务
	Schedule(spec string, jobId string, job cron.Job) (id cron.EntryID, err error)
	// Every 添加任务
	Every(duration time.Duration, job cron.Job)
	// Now 运行一个任务
	Now(job cron.Job)
	// In 定时运行一个任务
	In(duration time.Duration, job cron.Job)
	// Remove 移除一个任务
	Remove(jobId string)
	// Stop 功能停止运行
	Stop()
	// FindTask 任务是否存在
	FindTask(jobId string) (ok bool)
	// StatusJson 当前任务状态
	StatusJson() map[string]interface{}
}
```

- ### 定时任务提供方法

模块提供了对应的方法进行定时任务的管理

1. #### `Start` 定时任务初始化

方法定义：
```golang
// Start 初始化功能
// 参数 v：暂留
Start(v ...int)
```

使用示例：

在使用定时任务之前需要先进行初始化，定时任务模块在初始化时已经进行了初始化，所以在使用`OwnerTask`时无需关心初始化

```golang
func init() {
	OwnerTask = NewTask()
	OwnerTask.Start()
}
```

执行结果：
```
[JobRunner] 2022/09/23 - 10:51:26 Started... 
```

2. #### `Schedule` 添加定时任务

方法定义： 
```golang

// 参数：spec 时间表达式
// 参数：jobId 定时任务ID，唯一
// 参数：job 任务对象，需要实现 Run 方法
// 返回值：id 在定时任务中的 id
// 返回值：err 错误信息
Schedule(spec string, jobId string, job cron.Job) (id cron.EntryID, err error)
```

使用示例：

```golang
type Test struct {
	Name string
}

func (t Test) Run() {
	fmt.Println(t.Name, "任务执行", time.Now().Format("2006-01-02 15:04:05.000"))
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
```

执行结果：
```
[JobRunner] 2022/09/23 - 10:51:26 Started... 
=== RUN   TestSchedule
测试任务2-Schedule 任务执行 2022-09-23 10:51:28.003
测试任务-Schedule 任务执行 2022-09-23 10:51:28.003
测试任务-Schedule 任务执行 2022-09-23 10:51:30.003
测试任务2-Schedule 任务执行 2022-09-23 10:51:30.003
测试任务2-Schedule 任务执行 2022-09-23 10:51:32.002
测试任务-Schedule 任务执行 2022-09-23 10:51:32.002
测试任务-Schedule 任务执行 2022-09-23 10:51:34.003
测试任务2-Schedule 任务执行 2022-09-23 10:51:34.003
测试任务2-Schedule 任务执行 2022-09-23 10:51:36.003
测试任务-Schedule 任务执行 2022-09-23 10:51:36.003
--- PASS: TestSchedule (10.00s)
PASS
ok      granada.framework/server/task   10.212s


> 测试运行完成时间: 2022/9/23 10:51:36 <
```

3. #### `Every` 间隔时间执行任务

方法定义：

```golang
// Every 添加任务
// duration time.Duration 时间间隔
// 参数：jobId string 任务ID
// 参数：job cron.Job 任务对象
Every(duration time.Duration, jobId string, job cron.Job)
```

使用示例：

```golang
type Test struct {
	Name string
}

func (t Test) Run() {
	fmt.Println(t.Name, "任务执行", time.Now().Format("2006-01-02 15:04:05.000"))
}

func TestEvery(t *testing.T) {
	OwnerTask.Every(time.Duration(3*time.Second), "123", Test{Name: "测试任务-Every"})
	time.Sleep(10 * time.Second)
}
```

执行结果：

```
[JobRunner] 2022/09/23 - 11:21:25 Started... 
=== RUN   TestEvery
测试任务-Every 任务执行 2022-09-23 11:21:28.001
测试任务-Every 任务执行 2022-09-23 11:21:31.008
测试任务-Every 任务执行 2022-09-23 11:21:34.011
--- PASS: TestEvery (10.00s)
PASS
ok      granada.framework/server/task   10.238s


> 测试运行完成时间: 2022/9/23 11:21:35 <
```

4. #### `Now`立即执行任务

方法定义：

```golang
// Now 运行一个任务
// 参数：job cron.Job 任务对象
Now(job cron.Job)
```

使用示例：

```golang
type Test struct {
	Name string
}

func (t Test) Run() {
	fmt.Println(t.Name, "任务执行", time.Now().Format("2006-01-02 15:04:05.000"))
}

func TestRunNowTask(t *testing.T) {
	OwnerTask.Now(Test{Name: "测试任务马上执行"})
	time.Sleep(1 * time.Second)
}
```

执行结果：

```
[JobRunner] 2022/09/23 - 11:25:19 Started... 
=== RUN   TestRunNowTask
测试任务马上执行 任务执行 2022-09-23 11:25:19.821
--- PASS: TestRunNowTask (1.01s)
PASS
ok      granada.framework/server/task   1.230s


> 测试运行完成时间: 2022/9/23 11:25:20 <
```

5. #### `In`延迟运行任务一次

方法定义：

```golang
// In 延迟运行任务一次
// 参数：duration time.Duration 间隔时间
// 参数：job cron.Job 任务对象
In(duration time.Duration, job cron.Job)
```

使用示例：

```golang
type Test struct {
	Name string
}

func (t Test) Run() {
	fmt.Println(t.Name, "任务执行", time.Now().Format("2006-01-02 15:04:05.000"))
}

func TestIn(t *testing.T) {
	OwnerTask.In(time.Duration(3*time.Second), Test{Name: "测试任务-In"})
	time.Sleep(10 * time.Second)
}
```

执行结果：

```
[JobRunner] 2022/09/23 - 11:29:31 Started... 
=== RUN   TestIn
测试任务-In 任务执行 2022-09-23 11:29:34.283
--- PASS: TestIn (10.01s)
PASS
ok      granada.framework/server/task   10.290s


> 测试运行完成时间: 2022/9/23 11:29:41 <
```

6. #### `Remove` 移除一个任务

方法定义：

```golang

// Remove 移除一个任务
// 参数：jobId string 任务ID
Remove(jobId string)
```

使用示例：

```golang
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
```

执行结果：

```
[JobRunner] 2022/09/23 - 13:10:21 Started... 
=== RUN   TestRemove
测试任务-Remove 任务执行 2022-09-23 13:10:22.001
测试任务-Remove 任务执行 2022-09-23 13:10:24.014
测试任务-Remove 任务执行 2022-09-23 13:10:26.005
移除任务 2022-09-23 13:10:27.012
结束 2022-09-23 13:10:30.019
--- PASS: TestRemove (8.02s)
PASS
ok      granada.framework/server/task   8.232s


> 测试运行完成时间: 2022/9/23 13:10:30 <
```

7. #### `Stop`停止定时任务功能

方法定义：

```golang
// Stop 功能停止运行 注：在停用之后再次启用时所有任务失效
Stop()
```

使用示例：

```golang
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
```

执行结果：

```
[JobRunner] 2022/09/23 - 13:22:09 Started... 
=== RUN   TestStop
测试任务-Stop 任务执行 2022-09-23 13:22:10.014
测试任务-Stop 任务执行 2022-09-23 13:22:12.013
测试任务-Stop 任务执行 2022-09-23 13:22:14.003
停用定时任务功能 2022-09-23 13:22:14.272
启用定时任务功能 2022-09-23 13:22:19.287
[JobRunner] 2022/09/23 - 13:22:19 Started... 
--- PASS: TestStop (15.03s)
PASS
ok      granada.framework/server/task   15.248s


> 测试运行完成时间: 2022/9/23 13:22:24 <
```

8. #### `FindTask`查找定时任务是否存在

方法定义：

```golang
// FindTask 任务是否存在
// 参数：jobId string 任务ID
// 返回值： ok bool true 存在 false 不存在
FindTask(jobId string) (ok bool)
```

使用示例：

```golang
type Test struct {
	Name string
}

func (t Test) Run() {
	fmt.Println(t.Name, "任务执行", time.Now().Format("2006-01-02 15:04:05.000"))
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
```

执行结果：

```
[JobRunner] 2022/09/23 - 13:28:40 Started... 
=== RUN   TestFindTask
测试任务-Stop 任务执行 2022-09-23 13:28:42.005
测试任务-Stop 任务执行 2022-09-23 13:28:44.007
定时任务存在
测试任务-Stop 任务执行 2022-09-23 13:28:46.013
--- PASS: TestFindTask (5.01s)
PASS
ok      granada.framework/server/task   5.228s


> 测试运行完成时间: 2022/9/23 13:28:46 <
```

9. #### `StatusJson`当前定时任务状态

方法定义：

```golang
// StatusJson 当前任务状态
// 返回值：map[string]interface{} 当前定时任务状态
StatusJson() map[string]interface{}
```

使用示例：

```golang
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
```

执行结果：

```
[JobRunner] 2022/09/23 - 13:34:05 Started... 
=== RUN   TestStatusJson
测试任务-StatusJson 任务执行 2022-09-23 13:34:06.010
测试任务-StatusJson 任务执行 2022-09-23 13:34:08.000
测试任务-StatusJson 任务执行 2022-09-23 13:34:10.012
jsonData {"jobrunner":[{"Id":1,"JobRunner":{"Name":"测试任务-StatusJson","Status":"IDLE","Latency":"0s"},"Next":"2022-09-23 13:34:12","Prev":"2022-09-23 13:34:10"}]}
--- PASS: TestStatusJson (5.00s)
PASS
ok      granada.framework/server/task   5.245s


> 测试运行完成时间: 2022/9/23 13:34:10 <
```