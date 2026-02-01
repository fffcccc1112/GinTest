package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"test/config"
)

// 公共日志工具
var logger *zap.Logger

func Init(cfg *config.Config) {
	//设置日志级别
	level, err := zapcore.ParseLevel(cfg.LoggerConfig.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	//编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	//输出方式
	writer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(getFileWriter(cfg.LoggerConfig.Filename)),
	)
	//创建logger
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writer,
		level)
	logger = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
}

// getFileWriter 获取文件写入器（自动创建目录）
func getFileWriter(filename string) *os.File {
	// 创建日志目录
	dir := "./logs"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.Mkdir(dir, 0755)
	}

	// 打开文件（追加模式）
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logger.Error("打开日志文件失败", zap.Error(err))
		return os.Stdout
	}
	return file
}

// 封装日志方法
// 可变参数，变长参数语法糖
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}
