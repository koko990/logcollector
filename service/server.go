package service

import (
	"context"
	"github.com/koko990/logcollector/service/KuberHealth"
	"github.com/koko990/logcollector/service/source"
	"github.com/koko990/logcollector/service/storage"
	"github.com/koko990/logcollector/util"
	"strconv"
	"time"
)

const (
	defaultTryKuberTime = 10
	defaultTryKuberLink = 10
	_APIMinDuration     = 10
)

//
func (svr *GenericCollectorServer) taskStart(ctx context.Context, cancel context.CancelFunc) {
	var s = source.NewSourceFactory(svr.kuberAPIServerAddress)

	s.FetchResource()
}

//
func CollectorServerFactory(kuberAddress string) *GenericCollectorServer {
	return &GenericCollectorServer{kuberAPIServerAddress: kuberAddress}
}

//The entrance of cmd.
func (svr *GenericCollectorServer) RunServer(signChan chan struct{}) {
	svr.statusMonitor.NewStatus()
	util.Logger.SetInfo("main routine is run")
	svr.serverDaemon(signChan)

}

//
func (svr *GenericCollectorServer) checkHealth() {
	for i := 0; i < defaultTryKuberLink; i++ {
		KuberHealth.KuberMonitor(svr.kuberAPIServerAddress, svr.kuberMonitor)
		*svr.kuberMonitor = false
		time.Sleep(time.Second * defaultTryKuberTime)
	}

}

//
func (svr *GenericCollectorServer) loopTask(signChan chan struct{}) {
	ticker := time.NewTicker(time.Millisecond * 5000)
	for range ticker.C {
		storage.TimestampIndex = strconv.Itoa(int(time.Now().Unix()))
		switch *svr.restAPIMonitor && *svr.kuberMonitor {
		case true:
			util.Logger.SetInfo("logcollector is start")
			svr.taskRunWithContext()
			signChan <- struct{}{}
		case false:
			util.Logger.SetInfo("logcollector is Suspend,restAPIMonitor is :", *svr.restAPIMonitor,
				"kuberMonitor :", *svr.kuberMonitor)
		}
	}
}

//
func (svr *GenericCollectorServer) taskRunWithContext() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*4000)
	svr.taskStart(ctx, cancel)
	select {
	case <-ctx.Done():
		return
	}
}

//
func (svr *GenericCollectorServer) serverDaemon(signChan chan struct{}) {
	go svr.restTrigger()
	go svr.checkHealth()
	svr.loopTask(signChan)

}
