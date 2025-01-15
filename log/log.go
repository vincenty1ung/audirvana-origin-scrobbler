package log

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sys/unix"
	"gopkg.in/natefinch/lumberjack.v2"
)

var dayRotationCounter int32 = 0 // 记录每天日志文件计数值
var Logger *zap.Logger

// 创建日期分隔的日志文件
func openDailyLogFile(dir string) (zapcore.WriteSyncer, error) {
	var fsyncer zapcore.WriteSyncer
	atomic.AddInt32(&dayRotationCounter, 1)
	logFileName := time.Now().Format("2006-01-02") + ".log"

	filePath := filepath.Join(dir, logFileName)
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		fsyncer = zapcore.AddSync(io.Discard)
	}
	fsyncer = zapcore.AddSync(f)
	return fsyncer, nil
}

func createLumberJackLogger(filename string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename, // 文件位置
		MaxSize:    1,        // 进行切割之前,日志文件的最大大小(MB为单位)
		MaxAge:     7,        // 保留旧文件的最大天数
		MaxBackups: 10,       // 保留旧文件的最大个数
		Compress:   true,     // 是否压缩/归档旧文件
		LocalTime:  true,
	}
	// AddSync 将 io.Writer 转换为 WriteSyncer。
	// 它试图变得智能：如果 io.Writer 的具体类型实现了 WriteSyncer，我们将使用现有的 Sync 方法。
	// 如果没有，我们将添加一个无操作同步。

	return zapcore.AddSync(lumberJackLogger)
}

func LogInit(logPath, infoLevel string, c <-chan struct{}) *zap.Logger {
	err := os.MkdirAll(logPath, 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Error occurred:", err.Error())
		os.Exit(1)
	}
	atom := zap.NewAtomicLevel()
	switch infoLevel {
	case "debug":
		atom.SetLevel(zap.DebugLevel)
	case "info":
		atom.SetLevel(zap.InfoLevel)
	case "error":
		atom.SetLevel(zap.ErrorLevel)
	case "warn":
		atom.SetLevel(zap.WarnLevel)
	default:
		atom.SetLevel(zap.InfoLevel)

	}
	/*fsyncer, err := openDailyLogFile(logPath)
	if err != nil {
		fmt.Println("Failed to set up logger due to error:", err.Error())
		os.Exit(1)
	}*/
	fileAndStdoutSyncer := zapcore.NewMultiWriteSyncer(
		createLumberJackLogger(logPath+"/go_audirvana-origin-scrobbler.log"), zapcore.AddSync(os.Stdout),
	)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		fileAndStdoutSyncer,
		atom,
	)

	Logger = zap.New(
		core,
		// zap.Development(),
		zap.AddCaller(),
		zap.AddCallerSkip(0),
	)

	// 示例日志记录
	defer func() {
		err = Logger.Sync()
		if err != nil && (!errors.Is(err, unix.ENOTTY) || !errors.Is(err, syscall.ENOTTY) || !errors.Is(
			err, unix.EBADF,
		)) {
			// golang.org/x/sys/unix.ENOTTY (25)
			// golang.org/x/sys/unix.EBADF (9)
			fmt.Println(err)
		} else if err != nil {
			panic(err)
		} else {

		}
	}()

	// 定时任务逻辑（例如使用时间轮，cron定时任务等）
	// 这里仅示例代码中直接执行一次清理工
	go cleanOldLogs(logPath, c)
	return Logger
}

func cleanOldLogs(dir string, c <-chan struct{}) {
	ticker := time.NewTicker(time.Hour * 24)
	for {
		select {
		case <-c:
			fmt.Println("clean old logs exit")
			return
		case <-ticker.C:
			files, err := os.ReadDir(dir)
			if err != nil {
				fmt.Println("Failed to list log files:", err.Error())
				return
			}

			today := time.Now()
			for _, file := range files {
				info, _ := file.Info() // 获取FileInfo
				logTime := info.ModTime()

				if !file.IsDir() && logTime.Before(today.AddDate(0, 0, -7)) {
					err = os.Remove(filepath.Join(dir, file.Name()))
					if err != nil {
						fmt.Printf("Failed to remove old log: %s. Err: %v\n", filepath.Join(dir, file.Name()), err)
					} else {
						fmt.Println("Removed outdated log:", filepath.Join(dir, file.Name()))
					}
				}
			}
		}
	}

}
