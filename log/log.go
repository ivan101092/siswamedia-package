package log

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	zapLogger  *zap.Logger
	zapTracing *zap.Logger
)

func Init(cfg *Log) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed find current Directory for setup Log file | %s", err.Error())
	}

	fileLocation := currentDirectory + cfg.Location
	fileFormat := cfg.FileFormat

	rotateLog, _ := rotatelogs.New(
		fileLocation+fileFormat,
		rotatelogs.WithMaxAge(time.Duration(cfg.MaxAge)*24*time.Hour),          // Maximum time before deleting file log
		rotatelogs.WithRotationTime(time.Duration(cfg.RotationFile)*time.Hour), // Time before creating new file
		rotatelogs.WithClock(rotatelogs.Local),
		rotatelogs.WithLinkName(fileLocation+cfg.FileLinkName))

	encoder := zap.NewProductionEncoderConfig()
	encoder.TimeKey = "time"
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder

	logToFileEncoder := zapcore.NewCore(zapcore.NewJSONEncoder(encoder), zapcore.AddSync(rotateLog), zapcore.InfoLevel)
	consoleEncoder := zapcore.NewCore(zapcore.NewJSONEncoder(encoder), zapcore.AddSync(os.Stdout), zapcore.InfoLevel)

	core := zapcore.NewTee(logToFileEncoder)
	if cfg.LogToTerminal {
		core = zapcore.NewTee(logToFileEncoder, consoleEncoder)
	}

	if cfg.UseStackTrace {
		zapLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(1))
		zapTracing = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(3))
	} else {
		zapLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		zapTracing = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(3))
	}
}

func Info(msg string) {
	zapLogger.Info(msg)
}

func Infof(msg string, v ...interface{}) {
	zapLogger.Info(fmt.Sprintf(msg, v...))
}

func Warn(msg string) {
	zapLogger.Warn(msg)
}

func Warnf(ctx context.Context, msg string, v ...interface{}) {
	errMessage := fmt.Sprintf(msg, v...)
	appendErrorData(ctx, errMessage)
	zapLogger.Warn(errMessage, zap.String("ProcessID", cast.ToString(ctx.Value(ProcessIDContextKey))))
}

func Error(msg string) {
	zapLogger.Error(msg)
}

func Errorf(msg string, v ...interface{}) {
	zapLogger.Error(fmt.Sprintf(msg, v...))
}

func Fatal(msg string) {
	zapLogger.Fatal(msg)
}

func Fatalf(msg string, v ...interface{}) {
	zapLogger.Fatal(fmt.Sprintf(msg, v...))
}

// Special log only for tracking incoming request to Market Service
func RequestLog(data RequestLogModel) {
	fmt.Print("\033[32m")
	zapLogger.Info("REQUEST_LOG",
		zap.String("ProcessID", data.ProcessID),
		zap.String("IP", data.IP),
		zap.Int("UserID", data.UserID),
		zap.String("URL", data.URL),
		zap.String("HttpMethod", data.Method),
		zap.Any("RequestHeader", data.ReqHeader),
		zap.Any("RequestBody", data.ReqBody),
		zap.Any("ResponseHeader", data.RespHeader),
		zap.Any("ResponseBody", data.RespBody),
		zap.Any("Error", data.Error),
		zap.Int("StatusCode", data.StatusCode),
		zap.Int64("RequestDuration", data.Duration),
	)
	fmt.Print("\033[0m")
}

// Special log only for tracing service to service communication
func TracingLog(ctx context.Context, url, method string, resCode int, resPayload []byte, reqHeader, payload, respHeader interface{}, err error, dur int64) {
	var responsePayload interface{}
	if err := json.Unmarshal(resPayload, &responsePayload); err != nil {
		responsePayload = string(resPayload)
	}

	zapTracing.Info("TRACING_LOG",
		zap.String("ProcessID", cast.ToString(ctx.Value(ProcessIDContextKey))),
		zap.String("URL", url),
		zap.String("HttpMethod", method),
		zap.Any("RequestHeader", reqHeader),
		zap.Any("RequestBody", payload),
		zap.Any("ResponseHeader", respHeader),
		zap.Any("ResponseBody", responsePayload),
		zap.Int("StatusCode", resCode),
		zap.Int64("RequestDuration", dur),
		zap.Any("Error", err),
	)
}

// appendErrorData to TDR log
func appendErrorData(ctx context.Context, errMessage string) {
	if listError, ok := ctx.Value("errorList").(*Errors); ok {
		var caller, fileName string
		if _, file, line, ok := runtime.Caller(2); ok {
			if currentDirectory, err := os.Getwd(); err == nil {
				fileName = strings.TrimPrefix(file, currentDirectory)
			}
			caller = fmt.Sprintf("%s:%v", fileName, line)
		}

		*listError = append(*listError, ErrorData{
			Location: caller,
			Error:    errMessage,
		})
	}
}
