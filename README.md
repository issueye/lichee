
<div align="center">
<img src="resources/lichee.png" style="zoom: 30%;" />

# LICHEE ：一个运行 `javascript` 脚本的小工具

 支持`javascript` 脚本执行、定时任务、`http server`

</div>

演示视频：



- 定时获取数据库内容

<img src="resources/lichee.gif" style="zoom:80%;" />



- 爬虫网页内容

  <img src="resources/go-query.gif" style="zoom:70%;" />



## 配置文件

使用 `LICHEE` 需要先创建一个 `config.json` 的配置文件，在配置文件中添加中添加对应的配置信息

```json
{
  "local_port": 10066,
  "use_db": false,
  "log": {
    "path": "/lichee",
    "max_size": 10,
    "max_age": 10,
    "max_backups": 10,
    "compress": true,
    "level": -1
  },
  "db": {
    "user": "",
    "password": "",
    "host": "",
    "database": "",
    "port": 1433,
    "log_mode": true
  },
  "job": [
    {
      "name": "a test job",
      "id": 1,
      "expr": "0/5 * * * * ?",
      "benable": true,
      "path": "test.js"
    }
  ]
}
```



- `local_port` 提供服务的端口号
- `use_db` 是否使用数据库
- `log` 日志配置
  - `path` 日志输出路径
  - `max_size` 一个文件的最大大小
  - `max_age` 存放时间
  - `max_backups` 最大备份数据
  - `compress` 是否压缩
  - `level` 日志输出等级  注：请参照 `zap` 的日志等级
- `db` 数据库配置
  - `user` 用户名称
  - `password` 密码
  - `host` 服务地址
  - `database` 数据库
  - `port` 端口号
  - `log_mode` 是否输出 `SQL`
- `job` 任务配置
  - `name` 任务名称
  - `id` 任务ID
  - `expr` 时间表达式  注: 请参照 `CRON` 表达式
  - `benable` 是否启用
  - `path` 脚本路径

> 在 `use_db` 为 `false` 时不加载 `db` 数据库配置节点，数据库当前只添加了对 `sqlserver` 数据库的支持，后续会添加其他数据库的支持
>
> `job` 任务配置节点是一个数组，可配置多个定时任务，定时任务的时间表达式使用的是 `cron` 表达式，请自行了解 `cron` 相关的知识
>
> `LICHEE` 的定时任务的脚本路径默认在 `runtime/js` 文件夹下，在运行时会默认到此目录下查找脚本文件
>
> `log` 日志输出路径默认在 `runtime/logs` 文件夹下的对应日志输出路径下



## 提供的 `javascript`  `api` 方法

- `path/filepath` 
  - `abs`
  - `join`
  - `ext`
- `utils`
  - `print`
  - `panic`
  - `toString`
  - `toBase64`
  - `md5`
  - `sha1`
  - `arrayToMap`
- `types`
  - `newInt`
  - `intValue`
  - `newBool`
  - `boolValue`
  - `newString`
  - `stringValue`
  - `makeByteSlice`
  - `test`
  - `err`
  - `retUndefined`
  - `retNull`
- `time`
  - `sleep`
  - `nowString`
  - `nowDate`
  - `nowYear`
  - `nowMonth`
  - `nowDay`
  - `nowHour`
  - `nowMinute`
  - `nowSecond`
- `os`
  - `O_CREATE`  `O_WRONLY`  `O_RDONLY`  `O_RDWR`  `O_APPEND`  `O_EXCL`  `O_SYNC`  `O_TRUNC` 
  - `args`
  - `tempDir`
  - `hostname`
  - `getEnv`
  - `remove`
  - `removeAll`
  - `mkdir`
  - `mkdirAll`
  - `getwd`
  - `chdir`
  - `openFile`
  - `create`
  - `open`
  - `stat`
- `ini`
  - `create`
  - `getStr`
  - `getInt`
  - `getBool`
  - `getSection`
  - `setStr`
  - `setInt`
  - `setBool`
  - `save`
- `fmt`
  - `sprintf`
  - `printf`
  - `println`
  - `print`
- `file`
  - `write`
  - `read`
- `error`
  - `new`
- `db/local`
  - `query`
  - `exec`
  - `begin`
    - `commit`
    - `rollback`
    - `exec`
    - `query`
- `os/exec`
  - `command`
