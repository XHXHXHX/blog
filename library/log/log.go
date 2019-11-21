package log

import (
	"blog/library"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"sync"
	"time"
)

const FilePath = "../../logs"

type Hook struct {
	// Entries is an array of all entries that have been received by this hook.
	// For safe access, use the AllEntries() method, rather than reading this
	// value directly.
	Entries []logrus.Entry
	mu      sync.RWMutex
}

type MyLog struct {
	logger *logrus.Logger
	hook *Hook
}

var myLog *MyLog

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.WarnLevel)

	logger := logrus.New()
	hook := new(Hook)
	logger.Hooks.Add(hook)

	filename := logFileName()
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"filename" : filename,
		}).Error("Open log failed")
		return
	}
	logger.SetOutput(file)

	myLog = &MyLog{
		logger: logger,
		hook: hook,
	}
}

func New() *logrus.Logger {
	return myLog.logger
}

func Close() {
	myLog.logger.Exit(0)
}

func (t *Hook) Fire(e *logrus.Entry) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Entries = append(t.Entries, *e)
	return nil
}

func (t *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func logFileName() string {
	today := currentDate()
	sys_path := library.GetCurrentPath()
	today = today + ".log"

	return path.Join(sys_path, FilePath, today)
}

func currentDate() string {
	return time.Now().Format("20060102")
}