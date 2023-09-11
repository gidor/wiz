package logwrapper

import (
	"io"
	"log"
	"os"
)

var (
	flogger *log.Logger
	// dlogger *log.Logger
	ok bool
)

// func init() {
// 	dlogger = log.New(os.Stderr, "", log.LstdFlags)
// }

func SetOutput(w io.Writer) {
	flogger = log.New(w, "wiz", log.LstdFlags)
	ok = true
}

func Print(v ...any) {
	log.Print(v...)
	if ok {
		flogger.Print(v...)
	}
}
func Println(v ...any) {
	log.Println(v...)
	if ok {
		flogger.Println(v...)
	}
}

func Printf(format string, v ...any) {
	log.Printf(format, v...)
	if ok {
		flogger.Printf(format, v...)
	}
}

func Fatalf(format string, v ...any) {
	Printf(format, v...)
	os.Exit(1)
}
func Fatal(v ...any) {
	Print(v...)
	os.Exit(1)
}

func Running(task string, params map[string]string) {
	msg := "Run" + task
	for k, v := range params {
		msg += " " + k + "=" + v
	}
	Println(msg)
}
