package status

import "sync"

var statusCache *bool
var requestQueue *sync.Cond

func init()  {
	requestQueue=sync.NewCond()
}
func setEnable(m *sync.Mutex, statusCache *bool) {
	m.Lock()
	*statusCache = true
	m.Unlock()

}
func setUnable(m *sync.Mutex, statusCache *bool) {
	m.Lock()
	*statusCache = false
	m.Unlock()
	sync.NewCond()
}

type statusMonitor struct {
	kuberMonitor   *bool
	restAPIMonitor *bool
	restAPIMutex   *sync.Mutex
	enableRestCh   chan struct{}
	disableRestCh  chan struct{}
}
