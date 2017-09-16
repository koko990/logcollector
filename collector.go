package main

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/koko990/logcollector/cmd/app"
	"github.com/koko990/logcollector/util"
)

func main() {
	util.Logger.SetInfo("The cpu core is", runtime.NumCPU(), ",The app would use all of cores")
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())
	for {
		app.CollectorDaemon()
		time.Sleep(3 * time.Second)
	}

}
