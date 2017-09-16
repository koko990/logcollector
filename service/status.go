package service

import "time"

//
func (svr *statusMonitor) restTrigger() {
	select {
	case <-svr.enableRestCh:
		svr.updateStatus(true)
	case <-svr.disableRestCh:
		svr.updateStatus(false)
	}
}

//
func (svr *statusMonitor) updateStatus(status bool) {
	svr.restAPIMutex.Lock()
	*svr.restAPIMonitor = status
	svr.restAPIMutex.Unlock()
	time.Sleep(_APIMinDuration * time.Second)
}
