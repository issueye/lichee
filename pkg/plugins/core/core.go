package core

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	js "github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	_ "github.com/issueye/lichee/pkg/plugins/core/boltdb"
	_ "github.com/issueye/lichee/pkg/plugins/core/db"
	_ "github.com/issueye/lichee/pkg/plugins/core/goquery"
	_ "github.com/issueye/lichee/pkg/plugins/core/net/http"
	_ "github.com/issueye/lichee/pkg/plugins/core/net/url"
	_ "github.com/issueye/lichee/pkg/plugins/core/path"
	"go.uber.org/zap"
)

//go:embed js/*
var Script embed.FS

const (
	GoPlugins = "lichee"
)

var (
	globalConvertProg *js.Program                   // convert 转换代码的对应编译对象
	globalDayjsProg   *js.Program                   // dayjs 转换代码的对应编译对象
	LogMap            = make(map[string]*ZapLogger) // 日志对象
)

type ZapLogger struct {
	log   *zap.Logger
	Close func()
}

type ModuleFunc = func(vm *js.Runtime, module *js.Object)

// Core
// js运行时核心的结构体
type Core struct {
	// 全局js加载目录
	globalPath string
	// 外部添加到内部的内容
	pkg map[string]map[string]any
	// 外部注册的模块
	modules map[string]ModuleFunc
	// 编译之后的对象
	// pro *js.Program
	// 对应文件的编译对象
	proMap map[string]*js.Program
	// 日志对象
	logger *zap.Logger
	// 锁
	lock *sync.Mutex
	// Name
	name string
}

type OptFunc = func(*Core)

// NewCore
// 创建一个对象
func NewCore(opts ...OptFunc) *Core {
	c := new(Core)
	c.lock = new(sync.Mutex)
	c.pkg = make(map[string]map[string]any)
	// 初始化全局
	c.pkg[GoPlugins] = make(map[string]any)
	c.modules = make(map[string]func(vm *js.Runtime, module *js.Object))
	c.proMap = make(map[string]*js.Program)
	// 配置
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// OptionLog
// 配置日志
func OptionLog(log *zap.Logger) OptFunc {
	return func(core *Core) {
		core.logger = log
	}
}

func (c *Core) setupJSRuntime(rt *js.Runtime, logger *zap.Logger) error {
	// 输出日志
	console := newConsole(logger)
	o := rt.NewObject()
	o.Set("log", console.Log)
	o.Set("debug", console.Debug)
	o.Set("info", console.Info)
	o.Set("error", console.Error)
	o.Set("warn", console.Warn)

	err := rt.Set("console", o)
	if err != nil {
		return err
	}

	// 加载js模块
	c.loadScript(rt, "utils-arr2map", "convert.js", globalConvertProg)
	c.loadScript(rt, "dayjs", "dayjs.min.js", globalDayjsProg)

	return nil
}

func (c *Core) LoadModule(vm *js.Runtime) {
	// require.WithGlobalFolders(path)

	// 添加 导入方法 require
	registry := require.NewRegistry(
		// 全局加载路径
		require.WithGlobalFolders(c.globalPath),
	)
	registry.Enable(vm)

	// 添加 日志方法 console
	if c.name == "" {
		c.name = "lichee-test"
	}

	log, ok := LogMap[c.name]
	if !ok {
		l, close, err := newZap(fmt.Sprintf("runtime/logs/%s.log", c.name))
		if err != nil {
			c.Errorf("加载日志失败，失败原因：%s", err.Error())
		}
		log = &ZapLogger{
			log:   l,
			Close: close,
		}
		LogMap[c.name] = log
	}

	// 设置运行时
	c.setupJSRuntime(vm, log.log)

	// 加载全局对象
	c.loadVariable(vm)

	// 加载外部模块
	c.registerModule()
}

// GetRts
// 获取运行时
func (c *Core) GetRts() *js.Runtime {
	return js.New()
}

func (c *Core) SetGlobalPath(path string) {
	c.globalPath = path
}

// loadScript
// 加载文件中的js脚本
func (c *Core) loadScript(vm *js.Runtime, name string, jsName string, p *js.Program) {
	if p == nil {
		path := fmt.Sprintf("js/%s", jsName)
		src, err := Script.ReadFile(path)
		if err != nil {
			c.Errorf("读取文件失败，失败原因：%s", err.Error())
			return
		}
		p, err = js.Compile(name, string(src), false)
		if err != nil {
			return
		}
	}
	// 运行脚本
	_, err := vm.RunProgram(p)
	if err != nil {
		c.Errorf("运行脚本[%s]失败，失败原因：%s", name, err.Error())
	}
}

// Run
// 运行脚本 文件
func (c *Core) Run(name, path string) error {
	vm := c.GetRts()
	c.name = name
	return c.run(path, vm)
}

// RunVM
// 运行脚本 文件
func (c *Core) RunVM(path string, vm *js.Runtime) error {
	return c.run(path, vm)
}

func (c *Core) run(path string, vm *js.Runtime) error {
	c.LoadModule(vm)
	var tmpPath string
	if c.globalPath != "" {
		tmpPath = filepath.Join(c.globalPath, path)
	} else {
		tmpPath = path
	}

	// 读取文件
	src, err := os.ReadFile(tmpPath)
	if err != nil {
		c.Errorf("读取文件失败，失败原因：%s", err.Error())
	} else {
		// 编译文件
		pro, err := js.Compile(fmt.Sprintf("script_%s", time.Now().Format("20060102150405999")), string(src), false)
		if err != nil {
			c.Errorf("编译代码失败，失败原因：%s", err.Error())
		} else {
			c.proMap[path] = pro
		}
	}

	// 只有存在编译对象时，才运行
	if c.proMap[path] != nil {
		_, err := vm.RunProgram(c.proMap[path])
		if jsErr, ok := err.(*js.Exception); ok {
			return fmt.Errorf("运行脚本失败，失败原因：%s", jsErr.Error())
		}
	}
	return nil
}

// RunString
// 运行脚本 字符串
func (c *Core) RunString(src string) error {
	vm := c.GetRts()
	c.LoadModule(vm)
	_, err := vm.RunString(src)
	if jsErr, ok := err.(*js.Exception); ok {
		return fmt.Errorf("运行脚本失败，失败原因：%s", jsErr.Error())
	}
	return nil
}

// SetGlobalProperty
// 写入数据到全局对象中
func (c *Core) SetGlobalProperty(key string, value any) {
	// 添加锁
	c.lock.Lock()
	defer c.lock.Unlock()

	c.pkg[GoPlugins][key] = value
}

func (c *Core) loadVariable(vm *js.Runtime) {
	// 添加锁
	c.lock.Lock()
	defer c.lock.Unlock()

	// 加载其他模块
	for name, mod := range c.pkg {
		jsMod := vm.NewObject()
		for k, v := range mod {
			jsMod.Set(k, v)
		}
		vm.Set(name, jsMod)
	}
}

// registerModule
// 外部注册模块到js
func (c *Core) registerModule() {
	for Name, moduleFn := range c.modules {
		require.RegisterNativeModule(Name, func(runtime *js.Runtime, module *js.Object) {
			m := module.Get("exports").(*js.Object)
			moduleFn(runtime, m)
		})
	}
}

// SetProperty
// 向模块写入变量或者写入方法
func (c *Core) SetProperty(moduleName, key string, value any) {
	// 添加锁
	c.lock.Lock()
	defer c.lock.Unlock()

	mod, ok := c.pkg[moduleName]
	if !ok {
		c.pkg[moduleName] = make(map[string]any)
		mod = c.pkg[moduleName]
	}
	mod[key] = value
}

// RegisterModule
// 注册模块
func (c *Core) RegisterModule(moduleName string, fn ModuleFunc) {
	c.modules[moduleName] = fn
}
