package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	Log *zap.SugaredLogger
}

type LoggerConfig struct {
	OutputMode string // "terminal", "file", "both"
	LogLevel   string // "debug", "info", "warn", "error"
	LogDir     string // directory for log files
}

func NewLogger() Logger {
	return NewLoggerWithConfig(LoggerConfig{
		OutputMode: "both",   // default: both terminal and file
		LogLevel:   "debug",  // default: debug level
		LogDir:     "logger", // default: logger directory
	})
}

func NewLoggerWithConfig(config LoggerConfig) Logger {
	return Logger{
		Log: InitLogWithConfig(config),
	}
}

func Prefix() string {
	return "logger-" + time.Now().Format("2006-01-02")
}

func InitLog() *zap.SugaredLogger {
	return InitLogWithConfig(LoggerConfig{
		OutputMode: "both",
		LogLevel:   "debug",
		LogDir:     "logger",
	})
}

func InitLogWithConfig(config LoggerConfig) *zap.SugaredLogger {
	var cores []zapcore.Core
	encoder := getEncoder()
	level := getLogLevel(config.LogLevel)

	// Add terminal output if needed
	if config.OutputMode == "terminal" || config.OutputMode == "both" {
		terminalCore := zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), level)
		cores = append(cores, terminalCore)
	}

	// Add file output if needed
	if config.OutputMode == "file" || config.OutputMode == "both" {
		fileCore := zapcore.NewCore(encoder, getLogWriter(config.LogDir), level)
		cores = append(cores, fileCore)
	}

	// If no valid output mode, default to terminal
	if len(cores) == 0 {
		terminalCore := zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), level)
		cores = append(cores, terminalCore)
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.Development())
	sugarLogger := logger.Sugar()
	return sugarLogger
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.DebugLevel
	}
}

func getEncoder() zapcore.Encoder {
	loggerConfig := zap.NewProductionEncoderConfig()
	loggerConfig.TimeKey = "timestamp"
	loggerConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000Z07:00")
	loggerConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	loggerConfig.FunctionKey = "func"
	return zapcore.NewJSONEncoder(loggerConfig)
}

func getLogWriter(logDir string) zapcore.WriteSyncer {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		// If can't create directory, fallback to current directory
		logDir = "."
	}

	logFile := logDir + "/" + Prefix() + ".log"
	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	})
	return ws
}

// only message
func (Logs Logger) Info(msg string)    { Logs.Log.Info(msg) }
func (Logs Logger) Warning(msg string) { Logs.Log.Warn(msg) }
func (Logs Logger) Error(msg string)   { Logs.Log.Error(msg) }
func (Logs Logger) Fatal(msg string)   { Logs.Log.Fatal(msg) }
func (Logs Logger) Panic(msg string)   { Logs.Log.Panic(msg) }

// with data
func (Logs Logger) InfoW(msg string, data ...any)    { Logs.Log.Infow(msg, "data", data) }
func (Logs Logger) WarningW(msg string, data ...any) { Logs.Log.Warnw(msg, "data", data) }
func (Logs Logger) ErrorW(msg string, data ...any)   { Logs.Log.Errorw(msg, "data", data) }
func (Logs Logger) FatalW(msg string, data ...any)   { Logs.Log.Fatalw(msg, "data", data) }
func (Logs Logger) PanicW(msg string, data ...any)   { Logs.Log.Panicw(msg, "data", data) }

// with data and request id
func (Logs Logger) InfoT(msg, requestID string, data ...any) {
	Logs.Log.Infow(msg, "request-id", requestID, "data", data)
}
func (Logs Logger) WarningT(msg, requestID string, data ...any) {
	Logs.Log.Warnw(msg, "request-id", requestID, "data", data)
}
func (Logs Logger) ErrorT(msg, requestID string, data ...any) {
	Logs.Log.Errorw(msg, "request-id", requestID, "data", data)
}
func (Logs Logger) FatalT(msg, requestID string, data ...any) {
	Logs.Log.Fatalw(msg, "request-id", requestID, "data", data)
}
func (Logs Logger) PanicT(msg, requestID string, data ...any) {
	Logs.Log.Panicw(msg, "request-id", requestID, "data", data)
}

// Close logger and sync all buffered logs
func (Logs Logger) Close() {
	Logs.Log.Sync()
}
