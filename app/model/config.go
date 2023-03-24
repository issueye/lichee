package model

type Config struct {
	LocalPort int64        `json:"local_port"` // 本地端口号
	Log       *LogConfig   `json:"log"`        // 日志
	Db        *DbConfig    `json:"db"`         // 数据库
	Job       []*JobConfig `json:"job"`        // 定时任务
}

type LogConfig struct {
	Path       string `json:"path"`
	MaxSize    int    `json:"max_size"`    //文件大小限制,单位MB
	MaxAge     int    `json:"max_age"`     //日志文件保留天数
	MaxBackups int    `json:"max_backups"` //最大保留日志文件数量
	Compress   bool   `json:"compress"`    //是否压缩处理
	Level      int    `json:"level"`       // 等级
}

type DbConfig struct {
	Username string `json:"user"`     // 用户名称
	Password string `json:"password"` // 密码
	Host     string `json:"host"`     // 服务器地址
	Database string `json:"database"` // 数据库
	Port     int    `json:"port"`     // 端口号
	LogMode  bool   `json:"log_mode"` // 日志模式
}

type JobConfig struct {
	Name    string `json:"name"`    // 任务名称
	Id      int64  `json:"id"`      // 任务ID
	Expr    string `json:"expr"`    // 时间表达式
	Benable bool   `json:"benable"` // 是否启用
	Path    string `json:"path"`    // 任务路径
}
