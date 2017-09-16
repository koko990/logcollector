package service

import "sync"

type GenericCollectorServer struct {
	collectorLocalPort    string
	kuberAPIServerAddress string
	nodeCollector
	serviceCollector
	statusMonitor
}
type statusMonitor struct {
	kuberMonitor   *bool
	restAPIMonitor *bool
	restAPIMutex   *sync.Mutex
	enableRestCh   chan struct{}
	disableRestCh  chan struct{}
}

func (s *statusMonitor) NewStatus() {
	s.restAPIMonitor = new(bool)
	s.kuberMonitor = new(bool)
	s.restAPIMutex = new(sync.Mutex)
	s.enableRestCh = make(chan struct{})
	s.disableRestCh = make(chan struct{})
	*s.restAPIMonitor = true
	*s.kuberMonitor = false
}

type nodeCollector struct {
}

type serviceCollector struct {
}
