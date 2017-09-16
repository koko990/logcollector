package logger

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func (a *NewLog) newLog() {
	switch a.loggerInit.LogLevel {
	case LevelFatal:
		*a = NewLog{logInit{LevelDebug}, noneOut,
			noneOut, noneOut, noneOut, printInfoOut}
	case LevelError:
		*a = NewLog{logInit{LevelInfo}, noneOut,
			noneOut, noneOut, printErrOut, printInfoOut}
	case LevelWarn:
		*a = NewLog{logInit{LevelWarn}, noneOut,
			noneOut, printInfoOut, printErrOut, printInfoOut}
	case LevelInfo:
		*a = NewLog{logInit{LevelError}, noneOut,
			printInfoOut, printInfoOut, printErrOut, printInfoOut}
	case LevelDebug:
		*a = NewLog{logInit{LevelFatal}, printInfoOut,
			printInfoOut, printInfoOut, printErrOut, printInfoOut}
	}
}

func (a *NewLog) LogRegister(c LogLevel) {

	a.loggerInit.LogLevel = c
	a.newLog()

}

func noneOut(none ...interface{}) {
}
func info(args string) {
	fmt.Printf("%c[1;36m[INFO]%s %s \t%c[0m\n",
		0x1B, time.Now().Format("2006/01/02 - 15:04:05"), args, 0x1B)
}
func printInfoOut(parameter ...interface{}) {
	info(fmt.Sprint(parameter...))

}
func printErrOut(parameter ...interface{}) {
	errorLog(fmt.Sprint(parameter...))

}
func errorLog(args string) {
	fmt.Printf("%c[1;31m[ERROR]%s %s \t%c[0m\n",
		0x1B, time.Now().Format("2006/01/02 - 15:04:05"), args, 0x1B)
}
func (NewLog) HttpLog(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
