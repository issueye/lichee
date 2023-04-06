package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Path       string
	MaxSize    int  //文件大小限制,单位MB
	MaxAge     int  //日志文件保留天数
	MaxBackups int  //最大保留日志文件数量
	Compress   bool //是否压缩处理
	Level      int  // 等级
}

// InitLogger
/**
 * 初始化日志
 * filename 日志文件路径
 * level 日志级别
 * maxSize 每个日志文件保存的最大尺寸 单位: M
 * maxBackups 日志文件最多保存多少个备份
 * maxAge 文件最多保存多少天
 * compress 是否压缩
 * serviceName 服务名
 * 由于zap不具备日志切割功能, 这里使用lumberjack配合
 */
func InitLogger(cfg *Config) (*zap.SugaredLogger, *zap.Logger) {
	now := time.Now()
	LogsPath := cfg.Path
	infoLogFileName := fmt.Sprintf("%s/info/%04d-%02d-%02d.log", LogsPath, now.Year(), now.Month(), now.Day())
	errorLogFileName := fmt.Sprintf("%s/error/%04d-%02d-%02d.log", LogsPath, now.Year(), now.Month(), now.Day())
	var coreArr []zapcore.Core

	// 获取编码器
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "name",
		CallerKey:     "file",
		FunctionKey:   "func",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		// 时间格式
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(fmt.Sprintf("[%s]", t.Format("2006\\01\\02 15:04:05")))
		},
		EncodeDuration: func(duration time.Duration, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendInt64(int64(duration) / 1000000)
		},
		// 代码单元
		EncodeCaller: func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
			data := strings.Split(caller.String(), "/")
			encoder.AppendString(fmt.Sprintf(`[@%s]`, data[len(data)-1]))
		},
		// 日志级别
		EncodeLevel: func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(fmt.Sprintf(`[%s]`, level.String()))
		},
		ConsoleSeparator: " ",
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	// encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapcore.Level(cfg.Level))

	// 日志级别
	highPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.ErrorLevel
	})
	//lowPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
	//	return level < zap.ErrorLevel && level >= zap.DebugLevel
	//})

	// 当yml配置中的等级大于Error时，lowPriority级别日志停止记录
	//if cfg.Level >= 2 {
	//	lowPriority = zap.LevelEnablerFunc(func(level zapcore.Level) bool {
	//		return false
	//	})
	//}

	// info文件writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   infoLogFileName, //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    cfg.MaxSize,     //文件大小限制,单位MB
		MaxAge:     cfg.MaxAge,      //日志文件保留天数
		MaxBackups: cfg.MaxBackups,  //最大保留日志文件数量
		LocalTime:  false,           //
		Compress:   cfg.Compress,    //是否压缩处理
	})
	// 第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	infoFileCore := zapcore.NewCore(encoder,
		zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)),
		atomicLevel)

	// error文件writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   errorLogFileName, //日志文件存放目录
		MaxSize:    cfg.MaxSize,      //文件大小限制,单位MB
		MaxAge:     cfg.MaxAge,       //日志文件保留天数
		MaxBackups: cfg.MaxBackups,   //最大保留日志文件数量
		LocalTime:  false,
		Compress:   cfg.Compress, //是否压缩处理
	})
	// 第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority)

	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)
	logger := zap.New(zapcore.NewTee(coreArr...), zap.AddCaller())
	Log := logger.Sugar()

	return Log, logger
}
