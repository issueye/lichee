# <img src="resources/lichee.png" style="zoom: 30%;" />

# LICHEE ：一个运行 `javascript` 脚本的小工具

通过 `LICHEE` 提供的 `javascript` 脚本执行、定时任务、`http`服务，提供了一个轻量简单的动态工具



演示视频：

<img src="resources/lichee.gif" style="zoom:80%;" />



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
