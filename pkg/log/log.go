package log

import (
	"context"
	"fmt"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Sugare export var
var Sugare *zap.SugaredLogger

// LumberJackLogger export var
var LumberJackLogger *lumberjack.Logger

var l *Logger

type Logger struct {
	*zap.Logger
	opts *Options
}

type Options struct {
	CtxKey string //通过 ctx 传递 hlog 对象
}

type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	PanicLevel LogLevel = "panic"
	FatalLevel LogLevel = "fatal"
)

type LogFormat string

const (
	JsonType   LogFormat = "json"   // json 结构化字符串
	CommonType LogFormat = "common" // 普通格式字符串
)

type LogConfig struct {
	MaxCount int    // 日志文件保存最大数
	MaxSize  int    // 日志单个文件最大保存大小，单位为M
	Compress bool   // 自导打 gzip包 默认false
	FilePath string // 日志文件输出路径
	Level    LogLevel
	Format   LogFormat // 日志格式
}

// GetLogger returns logger
func GetLogger() *Logger {
	if l == nil {
		fmt.Println("Please initialize the hlog service first")
		return nil
	}
	return l
}

func (l *Logger) GetCtx(ctx context.Context) *zap.Logger {
	log, ok := ctx.Value(l.opts.CtxKey).(*zap.Logger)
	if ok {
		return log
	}
	return l.Logger
}

func (l *Logger) AddCtx(ctx context.Context, field ...zap.Field) (context.Context, *zap.Logger) {
	log := l.With(field...)
	ctx = context.WithValue(ctx, l.opts.CtxKey, log)
	return ctx, log
}

// Init def
func Init(config *LogConfig) {
	if config == nil {
		config = &LogConfig{
			MaxCount: 30,
			MaxSize:  10,
			Compress: true,
			FilePath: "./log/server.log",
			Level:    InfoLevel,
			Format:   JsonType,
		}
	}
	levelMap := map[LogLevel]zapcore.Level{
		DebugLevel: zap.DebugLevel,
		InfoLevel:  zap.InfoLevel,
	}
	LumberJackLogger = &lumberjack.Logger{
		Filename: config.FilePath, // 日志输出文件
		MaxSize:  config.MaxSize,  // 日志最大保存10M
		// MaxBackups: 5,  // 就日志保留5个备份
		MaxAge:   config.MaxCount, // 最多保留30天日志 和MaxBackups参数配置1个就可以
		Compress: config.Compress, // 自导打 gzip包 默认false
	}

	writer := zapcore.AddSync(LumberJackLogger)

	// 格式相关的配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 修改时间戳的格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 日志级别使用大写显示

	// 日志内容默认格式为普通格式字符串，设置 json 后为 json 结构化格式
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	if config.Format == JsonType {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewCore(encoder, writer, levelMap[config.Level])
	l = &Logger{
		opts: &Options{
			CtxKey: "log_key",
		},
	}
	l.Logger = zap.New(core, zap.AddCaller()) // 增加caller信息
	Sugare = l.Logger.Sugar()

	Sugare.Infof("zap log init ok.")
}
