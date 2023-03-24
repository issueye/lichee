// Package task 定时任务相关包；
// 可以引用 config、utils 和基础包下的 logger 包中的内容，可以被 initialize 包和 server 下的其他子包引用
package task

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// OwnerTask 全局对象
var OwnerTask ICronTaskMana

type ICronTaskMana interface {
	// Start 初始化功能
	Start(v ...int)
	// Schedule 添加任务
	Schedule(spec string, jobId string, job cron.Job) (id cron.EntryID, err error)
	// Every 添加任务
	Every(duration time.Duration, jobId string, job cron.Job)
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

type CronTaskMana struct {
	// MainCron 作业调度程序单例实例.
	MainCron *cron.Cron
	// 存储当前还在调度中的任务
	Name map[string]cron.EntryID
	// 锁
	lock sync.Mutex
}

// DefaultJobPoolSize 默认运行任务数量
const DefaultJobPoolSize = 10

var (
	// workPermits 可以运行任务的最大数量.
	workPermits chan struct{}

	// selfConcurrent 是否允许单个作业与其自身同时运行?
	selfConcurrent bool

	//green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	reset   = string([]byte{27, 91, 48, 109})

	functions = []interface{}{makeWorkPermits, isSelfConcurrent}
)

// Func
// Callers can use jobs.Func to wrap a raw func.
// (Copying the type to this package makes it more visible)
//
// For example:
//
//	jobrunner.Schedule("cron.frequent", jobs.Func(myFunc))
type Func func()

func (r Func) Run() { r() }

var secondsParser = cron.NewParser(
	cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	//cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
)

type StatusData struct {
	Id        cron.EntryID
	JobRunner *Job
	Next      string
	Prev      string
}

// makeWorkPermits 可以运行任务的最大数量
func makeWorkPermits(bufferCapacity int) {
	if bufferCapacity <= 0 {
		workPermits = make(chan struct{}, DefaultJobPoolSize)
	} else {
		workPermits = make(chan struct{}, bufferCapacity)
	}
}

// isSelfConcurrent 是否允许单个作业与其自身同时运行
func isSelfConcurrent(cocnurrencyFlag int) {
	if cocnurrencyFlag <= 0 {
		selfConcurrent = false
	} else {
		selfConcurrent = true
	}
}

func NewTask() ICronTaskMana {
	task := new(CronTaskMana)
	task.MainCron = &cron.Cron{}
	task.Name = make(map[string]cron.EntryID)
	task.lock = sync.Mutex{}
	return task
}

// Schedule 按照时间规则运行一个任务
func (t *CronTaskMana) Schedule(spec string, jobId string, job cron.Job) (id cron.EntryID, err error) {
	var sched cron.Schedule
	sched, err = secondsParser.Parse(spec)
	if err != nil {
		return -1, err
	}
	id = t.MainCron.Schedule(sched, New(job))
	t.Name[jobId] = id
	return
}

// Every 按照一个时间范围定期运行
func (t *CronTaskMana) Every(duration time.Duration, jobId string, job cron.Job) {
	id := t.MainCron.Schedule(cron.Every(duration), New(job))
	t.Name[jobId] = id
}

// Now 立即执行一个任务一次
func (t *CronTaskMana) Now(job cron.Job) {
	go New(job).Run()
}

// In
// 延迟运行任务一次.
func (t *CronTaskMana) In(duration time.Duration, job cron.Job) {
	go func() {
		time.Sleep(duration)
		New(job).Run()
	}()
}

// Start 开启一个服务
func (t *CronTaskMana) Start(v ...int) {
	t.MainCron = cron.New(cron.WithSeconds())
	t.Name = make(map[string]cron.EntryID)
	for i, option := range v {
		functions[i].(func(int))(option)
	}

	t.MainCron.Start()

	fmt.Printf("%s[JobRunner] %v Started... %s \n",
		magenta, time.Now().Format("2006/01/02 - 15:04:05"), reset)
}

// Stop 停止服务
func (t *CronTaskMana) Stop() {
	go t.MainCron.Stop()
}

// Remove 移除一个任务
func (t *CronTaskMana) Remove(jobId string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.MainCron.Remove(t.Name[jobId])
	delete(t.Name, jobId)
}

// FindTask 查找是否存在任务
func (t *CronTaskMana) FindTask(jobId string) (ok bool) {
	_, ok = t.Name[jobId]
	return
}

// Entries
// Return detailed list of currently running recurring jobs
// to remove an entry, first retrieve the ID of entry
func (t *CronTaskMana) Entries() []cron.Entry {
	return t.MainCron.Entries()
}

func (t *CronTaskMana) StatusPage() []StatusData {

	ents := t.MainCron.Entries()
	Statuses := make([]StatusData, len(ents))
	for k, v := range ents {
		Statuses[k].Id = v.ID
		Statuses[k].JobRunner = t.addJob(v.Job)
		Statuses[k].Next = v.Next.Format("2006-01-02 15:04:05.999")
		Statuses[k].Prev = v.Prev.Format("2006-01-02 15:04:05.999")
	}
	return Statuses
}

// StatusJson 任务的当前状态
func (t *CronTaskMana) StatusJson() map[string]interface{} {

	return map[string]interface{}{
		"jobrunner": t.StatusPage(),
	}
}

func (t *CronTaskMana) addJob(job cron.Job) *Job {
	return job.(*Job)
}

func init() {
	OwnerTask = NewTask()
	OwnerTask.Start()
}
