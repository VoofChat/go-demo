package main

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var sugarLogger1 *zap.SugaredLogger

func InitZap1() {
	//	1.配置zapcore的编码器
	zapEncode := zapcore.EncoderConfig{
		MessageKey:     "msg",   //消息的字段名
		LevelKey:       "level", //调试等级的字段名
		TimeKey:        "time",  //时间
		NameKey:        "name",
		CallerKey:      "line",
		FunctionKey:    "method",
		StacktraceKey:  "Stacktrace",
		SkipLineEnding: false,
		LineEnding:     zapcore.DefaultLineEnding,     //输出的分割符
		//EncodeLevel:    zapcore.LowercaseLevelEncoder, //序列化字符串的大小写
		EncodeLevel:    zapcore.CapitalLevelEncoder, //序列化字符串的大小写
		//EncodeTime:          zapcore.ISO8601TimeEncoder,     //时间的编码格式
		EncodeTime:          EncodeTime,                     //时间自定义的
		EncodeDuration:      zapcore.SecondsDurationEncoder, //时间显示的位数
		EncodeCaller:        zapcore.ShortCallerEncoder,     //输出的运行文件路径长度
		EncodeName:          zapcore.FullNameEncoder,        //可选的
		NewReflectedEncoder: nil,
		ConsoleSeparator:    "", //控制台格式时，每个字段间的分割符,不配置默认即可
	}
	//2.日志分割器
	hook := &lumberjack.Logger{
		Filename:   "logs/log1.log", //日志文件路径
		MaxSize:    128,            // 每个日志文件保存的大小 单位:M
		MaxAge:     7,              // 文件最多保存多少天
		MaxBackups: 30,             // 日志文件最多保存多少个备份
		Compress:   false,          // 是否压缩
	}

	//	3.设置日志
	//是否输出到控制台，默认为false
	isConsole := false
	//日志的编码格式，分为 json和Console
	encodeMode := "json"
	//输出级别
	logLev := zap.NewAtomicLevel()
	logMode := "debug"
	switch logMode {
	case "debug":
		isConsole = true
		logLev.SetLevel(zapcore.DebugLevel)
	case "info":
		logLev.SetLevel(zapcore.InfoLevel)
	case "warn":
		logLev.SetLevel(zapcore.WarnLevel)
	case "errors":
		logLev.SetLevel(zapcore.ErrorLevel)
	default:
		logLev.SetLevel(zapcore.DebugLevel)
	}
	//	4.设置zap日志输出位置，使用数组的方式便于控制输出到多个位置
	writes := []zapcore.WriteSyncer{zapcore.AddSync(hook)}
	//如果需要同步输出到控制台
	if isConsole {
		writes = append(writes, zapcore.AddSync(os.Stdout))
	}
	//设置日志的编码格式json和Console
	var enc zapcore.Encoder
	if encodeMode == "json" {
		enc = zapcore.NewJSONEncoder(zapEncode)
	} else {
		enc = zapcore.NewConsoleEncoder(zapEncode)
	}
	//	5.通过传入的配置实例化一个core
	core := zapcore.NewCore(enc,
		zapcore.NewMultiWriteSyncer(writes...),
		logLev)

	//6.构造日志
	//设置为开发模式会记录panic
	development := zap.Development()
	//开启记录文件名和行号
	caller := zap.AddCaller()
	//caller := zap.WithCaller(true)
	//构造一个字段
	zap.Fields(zap.String("appName", "demozap"))
	//通过传入的配置实例化一个日志
	zapLogger := zap.New(core, development, caller)
	zapLogger.Info("初始化日志")
	sugarLogger1 = zapLogger.Sugar()
}

// EncodeTime 自定义时间输出编码器
func EncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("我的:" + "2006/01/02 - 15:04:05.000"))
}
