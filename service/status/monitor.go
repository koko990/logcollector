package status

import "sync"

func (s *statusMonitor) NewStatus() {
	s.restAPIMonitor = new(bool)
	s.kuberMonitor = new(bool)
	s.restAPIMutex = new(sync.Mutex)
	s.enableRestCh = make(chan struct{})
	s.disableRestCh = make(chan struct{})
	*s.restAPIMonitor = true
	*s.kuberMonitor = false
}
