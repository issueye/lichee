package global

// 排队叫号服务运行模式类型
type IServerModeType int

const (
	SMT_DEFAULT IServerModeType = iota // 默认模式 适用于V20 V21
	SMT_VIEW                           // 视图模式 适用于 云HIS 或三方视图模式
	SMT_API                            // 接口模式 适用于三方接口方式（注：暂未实现，作为保留）
)

// 数据库类型
type IDbType int

const (
	DB_SSMS   IDbType = iota // sqlserver
	DB_MYSQL                 // mysql
	DB_ORACLE                // oracle
)

const (
	VIEW_TBZDMZYS = "VTBZDMZYS"
	VIEW_TBZDMZKS = "VTBZDMZKS"
)

const (
	CONF_SYSTEM = "SYSTEM"
	CONF_LOG    = "LOG"
	CONF_DB     = "DB"
	CONF_LINE   = "LINE"
)
