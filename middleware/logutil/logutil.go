package logutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	logger "github.com/sirupsen/logrus"
	"github.com/wuchunfu/JobFlow/middleware/config"
	"github.com/wuchunfu/JobFlow/utils/fileUtils"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func InitLog(setting *config.Log) {
	file, logFilePath := LogFile(setting.FilePath, setting.FileName)

	mode := strings.Replace(strings.ToLower(setting.Mode), " ", "", -1)
	switch mode {
	case "console":
		logger.SetFormatter(&logger.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
		logger.SetOutput(os.Stdout)
	case "file":
		logger.SetOutput(file)
	case "console,file":
		// 日志输出到文件同时输出到控制台
		multiWriter := io.MultiWriter(os.Stdout, file)
		logger.SetFormatter(&logger.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
		logger.SetOutput(multiWriter)
	default:
		logger.SetOutput(file)
	}

	// 设置日志级别
	level := strings.Replace(strings.ToLower(setting.Level), " ", "", -1)
	switch level {
	// 如果日志级别不是debug就不要打印日志到控制台了
	case "debug":
		// 设置显示文件名和行号, 将函数名和行数放在日志里面
		logger.SetReportCaller(true)
		logger.SetLevel(logger.DebugLevel)
		logger.SetOutput(os.Stderr)
	case "info":
		logger.SetLevel(logger.InfoLevel)
	case "warn":
		logger.SetLevel(logger.WarnLevel)
	case "error":
		// 设置显示文件名和行号, 将函数名和行数放在日志里面
		logger.SetReportCaller(true)
		logger.SetLevel(logger.ErrorLevel)
	default:
		logger.SetLevel(logger.InfoLevel)
	}
	// 设置rotatelogs日志分割Hook
	logger.AddHook(NewLfsHook(logFilePath))
}

func LogFile(logDirPath string, logFileName string) (*os.File, string) {
	if !fileUtils.IsExistPath(logDirPath) {
		if err := os.MkdirAll(logDirPath, 0755); err != nil {
			logger.Errorf("create log directory failed: %v\n", err)
			os.Exit(-1)
		}
	}

	// 文件路径
	logFilePath := path.Join(logDirPath, logFileName)

	_, err := WriteFile(logFilePath)
	if err != nil {
		if fileUtils.IsExistPath(logFilePath) {
			if err := os.Remove(logFilePath); err != nil {
				logger.Warnf("File deletion failed: %v\n", err)
			}
		}
	}
	file, _ := WriteFile(logFilePath)

	return file, logFilePath
}

func WriteFile(filePath string) (*os.File, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		//logrus.Errorf("Open File err: %v\n", err)
		return nil, err
	}
	return file, nil
}

func NewLfsHook(filePath string) logger.Hook {
	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		filePath+".%Y%m%d.log",
		// 生成软链，指向最新日志文件, WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
		rotatelogs.WithLinkName(filePath),
		// WithMaxAge和WithRotationCount二者只能设置一个,
		// 设置文件清理前的最长保存时间, 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// 设置文件清理前最多保存的个数
		//rotatelogs.WithRotationCount(5),
		// 设置日志分割的时间,这里设置为一天分割一次, 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	if err != nil {
		logger.Errorf("config logger error: %v\n", err)
	}

	writeMap := lfshook.WriterMap{
		logger.DebugLevel: logWriter,
		logger.InfoLevel:  logWriter,
		logger.WarnLevel:  logWriter,
		logger.ErrorLevel: logWriter,
		logger.FatalLevel: logWriter,
		logger.PanicLevel: logWriter,
	}

	lfsHook := lfshook.NewHook(writeMap, &logger.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return lfsHook
}

//日志自定义格式
type LogFormatter struct{}

//格式详情
func (logFormat *LogFormatter) Format(entry *logger.Entry) ([]byte, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	var file string
	var line int
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		line = entry.Caller.Line
	}
	level := strings.ToUpper(entry.Level.String())
	content, _ := json.Marshal(entry.Data)
	msg := fmt.Sprintf("%s [%s] [GOID:%d] [%s:%d] #msg:%s #content:%v\n",
		timestamp,
		level,
		getGID(),
		file,
		line,
		entry.Message,
		content,
	)
	return []byte(msg), nil
}

// 获取当前协程id
func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
