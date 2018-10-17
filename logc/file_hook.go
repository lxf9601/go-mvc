// 第三方库的源文件
// 主要修复了第三方代码行号的记录错误的BUG
package logc

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gogap/logrus"
	"github.com/gogap/logrus/hooks/caller"
	file2 "github.com/gogap/logrus/hooks/file"
)

func NewHook(file string) (f *FileHook) {
	path := strings.Split(file, "/")
	if len(path) > 1 {
		exec.Command("mkdir", path[0]).Run()
	}
	w := file2.NewFileWriter()
	config := fmt.Sprintf(`{"filename":"%s","maxdays":7}`, file)
	w.Init(config)
	return &FileHook{w}
}

type FileHook struct {
	W file2.LoggerInterface
}

func (hook *FileHook) Fire(entry *logrus.Entry) (err error) {
	message, err := getMessage(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}
	switch entry.Level {
	case logrus.PanicLevel:
		fallthrough
	case logrus.FatalLevel:
		fallthrough
	case logrus.ErrorLevel:
		return hook.W.WriteMsg(fmt.Sprintf("[ERROR] %s", message), file2.LevelError)
	case logrus.WarnLevel:
		return hook.W.WriteMsg(fmt.Sprintf("[WARN] %s", message), file2.LevelWarn)
	case logrus.InfoLevel:
		return hook.W.WriteMsg(fmt.Sprintf("[INFO] %s", message), file2.LevelInfo)
	case logrus.DebugLevel:
		return hook.W.WriteMsg(fmt.Sprintf("[DEBUG] %s", message), file2.LevelDebug)
	default:
		return nil
	}
	return
}

func (hook *FileHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

func getMessage(entry *logrus.Entry) (message string, err error) {
	message = message + fmt.Sprintf("%s\n", entry.Message)
	for k, v := range entry.Data {
		if !strings.HasPrefix(k, "err_") {
			message = message + fmt.Sprintf("%v:%v\n", k, v)
		}
	}
	if full, ok := entry.Data["err_full"]; ok {
		message = message + fmt.Sprintf("%v", full)
	} else {
		file, lineNumber := caller.GetCaller(2, "logc/logc.go", "logrus/hooks.go",
			"logrus/entry.go", "logrus/logger.go", "logrus/exported.go", "asm_amd64.s")
		if file != "" {
			sep := fmt.Sprintf("%s/src/", os.Getenv("GOPATH"))
			fileName := strings.Split(file, sep)
			if len(fileName) >= 2 {
				file = fileName[1]
			}
		}
		message = message + fmt.Sprintf("%s:%d", file, lineNumber)
	}
	return
}
