package logs

import (
	"fmt"
	"os"
	"time"

	"GoTenancy/files"
	"github.com/fatih/color"
)

func NewLog() *os.File {
	path := "./logs/"
	_ = files.CreateFile(path)
	filename := getFileName(path)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		color.Red(fmt.Sprintf("日志记录出错: %v", err))
	}
	return f
}

func getFileName(path string) string {
	return path + time.Now().Format("2006-01-02") + ".log"
}
