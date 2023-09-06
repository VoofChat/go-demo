package main

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var sugarLogger2 *zap.SugaredLogger

func InitZap2() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	//core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout)), zapcore.DebugLevel)
	// zap.AddCaller() 该配置可以增加 代码行号 看是哪个位置打印的
	logger := zap.New(core, zap.AddCaller())
	sugarLogger2 = logger.Sugar()
}

func getLogWriter() zapcore.WriteSyncer {
	//Zap本身不支持切割归档日志文件 使用第三方库 Lumberjack进行切割
	lumberjackLogger := &lumberjack.Logger{
		// 文件路径
		Filename: "logs/log2.log",
		// 日志最大保存
		MaxSize: 1,
		// 文件保留最大数量
		MaxBackups: 5,
		// 最大保留天数
		MaxAge:   30,
		Compress: false,
		// 使用本地时间命名 默认为utc时间命名
		LocalTime: true,
	}
	//file, _ := os.Create("./test.zap")
	return zapcore.AddSync(lumberjackLogger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	// 将TimeKey变为time, 默认为ts
	encoderConfig.TimeKey = "time"
	// 将时间戳变更为 2006-01-02 15:04:05
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 将输出等级转为大写 info -> INFO
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	//return zapcore.NewConsoleEncoder(encoderConfig)
	return zapcore.NewJSONEncoder(encoderConfig)
}


