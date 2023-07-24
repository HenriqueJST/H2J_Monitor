package logger

import (
	"fmt"

	"github.com/golang/glog"
)

func Info(args ...any) {
	fmt.Println(args...)
	glog.Info(args)
}

func Error(args ...any) {
	fmt.Println(args...)
	glog.Error(args)
}

func Infof(format string, args ...any) {
	fmt.Printf(format, args...)
	glog.Infof(format, args)
}

func Errorf(format string, args ...any) {
	fmt.Printf(format, args...)
	glog.Errorf(format, args)
}

func Debug(args ...any) {
	fmt.Println(args...)
}
