package core

import (
	"encoding/json"
	"strings"

	"github.com/dop251/goja"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type console struct {
	logger *zap.Logger
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 修改时间编码器

	// 在日志文件中使用大写字母记录日志级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// NewConsoleEncoder 打印更符合人们观察的方式
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(path string) (zapcore.WriteSyncer, func(), error) {
	return zap.Open(path)
}

func newZap(path string) (*zap.Logger, func(), error) {
	encode := getEncoder()
	ws, close, err := getLogWriter(path)
	if err != nil {
		return nil, nil, err
	}

	core := zapcore.NewCore(encode, ws, zap.DebugLevel)
	log := zap.New(core, zap.AddCaller())
	return log, close, nil
}

func newConsole(log *zap.Logger) *console {
	return &console{logger: log}
}

func (c console) log(level zapcore.Level, args ...goja.Value) {
	var strs strings.Builder
	for i := 0; i < len(args); i++ {
		if i > 0 {
			strs.WriteString(" ")
		}
		strs.WriteString(c.valueString(args[i]))
	}
	msg := strs.String()

	switch level { //nolint:exhaustive
	case zapcore.DebugLevel:
		c.logger.Debug(msg)
	case zapcore.InfoLevel:
		c.logger.Info(msg)
	case zapcore.WarnLevel:
		c.logger.Warn(msg)
	case zapcore.ErrorLevel:
		c.logger.Error(msg)
	}
}

func (c console) Log(args ...goja.Value) {
	c.Info(args...)
}

func (c console) Debug(args ...goja.Value) {
	c.log(zapcore.DebugLevel, args...)
}

func (c console) Info(args ...goja.Value) {
	c.log(zapcore.InfoLevel, args...)
}

func (c console) Warn(args ...goja.Value) {
	c.log(zapcore.WarnLevel, args...)
}

func (c console) Error(args ...goja.Value) {
	c.log(zapcore.ErrorLevel, args...)
}

func (c console) valueString(v goja.Value) string {
	mv, ok := v.(json.Marshaler)
	if !ok {
		return v.String()
	}

	b, err := json.Marshal(mv)
	if err != nil {
		return v.String()
	}
	return string(b)
}
