package log

import (
	"BeeScan-scan/pkg/config"
	"fmt"
	"github.com/fatih/color"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"reflect"
	"strings"
	"time"
)

/*
创建人员：云深不知处
创建时间：2022/1/13
程序功能：日志
*/
type Config interface {
	Level() zapcore.Level
	LogPath() string
	LogMaxSize() int
	InfoOutput() string
	ErrorOutput() string
	DebugOutput() string
}

func Setup() {

	LoggerConfig := config.GlobalConfig.LogConfig
	NewLogger(config.GlobalConfig.Level(), LoggerConfig.MaxSize, LoggerConfig.MaxBackups, LoggerConfig.MaxAge, LoggerConfig.Compress, config.GlobalConfig)
}

var (
	sugarInfoLogger  *zap.SugaredLogger
	sugarInfoPath    string
	sugarDebugLogger *zap.SugaredLogger
	sugarDebugPath   string
	sugarErrorLogger *zap.SugaredLogger
	sugarErrPath     string
)

func GetInfoLogPath() string {
	return sugarInfoPath
}

func GetDebugLogPath() string {
	return sugarDebugPath
}

func GetErrLogPath() string {
	return sugarErrPath
}

func formatArgs(v ...interface{}) string {
	var formatStrings []string
	for i := 0; i < len(v); i++ {
		t := v[i]
		switch reflect.TypeOf(t).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(t)
			for i := 0; i < s.Len(); i++ {
				formatStrings = append(formatStrings, `%v`)
			}
		}

	}
	return strings.Join(formatStrings, " ")
}

func Info(args ...interface{}) {
	format := formatArgs(args)
	sugarInfoLogger.Info("", fmt.Sprintf(format, args...))
	_, _ = fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", fmt.Sprintf(format, args...))
}

func Error(args ...interface{}) {
	format := formatArgs(args)
	sugarErrorLogger.Error("", fmt.Sprintf(format, args...))
	_, _ = fmt.Fprintln(color.Output, color.HiRedString("[ERRO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", fmt.Sprintf(format, args...))
}

func Warn(args ...interface{}) {
	format := formatArgs(args)
	sugarErrorLogger.Warn("", fmt.Sprintf(format, args...))
	_, _ = fmt.Fprintln(color.Output, color.HiYellowString("[WARN]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", fmt.Sprintf(format, args...))
}

func Debug(args ...interface{}) {
	format := formatArgs(args)
	sugarDebugLogger.Debug("", fmt.Sprintf(format, args...))
	_, _ = fmt.Fprintln(color.Output, color.HiGreenString("[DEBG]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", fmt.Sprintf(format, args...))
}

func Panic(args ...interface{}) {
	format := formatArgs(args)
	sugarErrorLogger.Panic("", fmt.Sprintf(format, args...))
	_, _ = fmt.Fprintln(color.Output, color.HiRedString("[FATA]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", fmt.Sprintf(format, args...))
}

func InfoOutput(args ...interface{}) {
	format := formatArgs(args)
	_, _ = fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", fmt.Sprintf(format, args...))
}

func WarningOutput(args ...interface{}) {
	format := formatArgs(args)
	_, _ = fmt.Fprintln(color.Output, color.HiYellowString("[WARN]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", fmt.Sprintf(format, args...))
}

func ErrorOutput(args ...interface{}) {
	format := formatArgs(args)
	_, _ = fmt.Fprintln(color.Output, color.HiRedString("[ERRO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", fmt.Sprintf(format, args...))
}

func DebugOutput(args ...interface{}) {
	format := formatArgs(args)
	_, _ = fmt.Fprintln(color.Output, color.HiGreenString("[DEBG]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", fmt.Sprintf(format, args...))
}

func VerboseOutput(args ...interface{}) {
	format := formatArgs(args)
	_, _ = fmt.Fprintln(color.Output, color.HiMagentaString("[VEBO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", fmt.Sprintf(format, args...))
}

func FatalOutput(args ...interface{}) {
	format := formatArgs(args)
	_, _ = fmt.Fprintln(color.Output, color.HiRedString("[FATA]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", fmt.Sprintf(format, args...))
}

func createLogger(path string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool) *zap.SugaredLogger {
	core := newCore(path, level, maxSize, maxBackups, maxAge, compress)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return logger.Sugar()
}

func NewLogger(level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool, config Config) {
	var infoPath, debugPath, errPath, logPath string

	if _, err := os.Stat(config.LogPath()); os.IsNotExist(err) {
		_ = os.Mkdir(config.LogPath(), 0755)
	}

	if len(config.LogPath()) == 0 {
		logPath = "BeeScanLogs"
	}
	infoPath = fmt.Sprintf("%s.log", logPath)
	if config.InfoOutput() != "" {
		infoPath = config.InfoOutput()
	}
	sugarInfoLogger = createLogger(infoPath, level, maxSize, maxBackups, maxAge, compress)
	sugarInfoPath = infoPath

	sugarDebugLogger = sugarInfoLogger
	sugarDebugPath = infoPath

	sugarErrorLogger = sugarInfoLogger
	sugarErrPath = infoPath

	if config.DebugOutput() != "" {
		debugPath = config.DebugOutput()
		sugarDebugLogger = createLogger(debugPath, level, maxSize, maxBackups, maxAge, compress)
		sugarDebugPath = debugPath
	}

	if config.ErrorOutput() != "" {
		errPath = config.ErrorOutput()
		sugarErrorLogger = createLogger(errPath, level, maxSize, maxBackups, maxAge, compress)
		sugarErrPath = errPath
	}

	// logger = zap.New(core, zap.AddCaller(), zap.Development(), zap.Fields(zap.String("serviceName", serviceName)))
}

/**
 * zapcore构造
 */
func newCore(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool) zapcore.Core {
	//日志文件路径配置2
	hook := lumberjack.Logger{
		Filename:   filePath,   // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   compress,   // 是否压缩
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	//公用编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 编码器配置
		//zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook)), // 打印到文件
		atomicLevel, // 日志级别
	)
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
