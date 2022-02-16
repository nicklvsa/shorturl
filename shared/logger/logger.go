package logger

import "fmt"

func Infof(data string, args ...interface{}) {
	out := fmt.Sprintf("[Info] - %s\n", data)
	fmt.Printf(out, args...)
}

func Warnf(data string, args ...interface{}) {
	out := fmt.Sprintf("[Warn] - %s\n", data)
	fmt.Printf(out, args...)
}

func Errorf(data string, args ...interface{}) {
	out := fmt.Sprintf("[Error] - %s\n", data)
	fmt.Printf(out, args...)
}

func Panicf(data string, args ...interface{}) {
	out := fmt.Sprintf("[Warn] - %s\n", data)
	panic(fmt.Sprintf(out, args...))
}
