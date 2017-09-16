package source

import (
	"encoding/json"

	"io/ioutil"
	"net/http"

	"github.com/koko990/logcollector/util"

	"github.com/koko990/logcollector/service/storage"

	"strings"

	"sync"

	"github.com/google/cadvisor/info/v2"
	modelK8s "k8s.io/api/core/v1"
)

var PodList modelK8s.PodList
var NodeList modelK8s.NodeList
var ServiceList modelK8s.ServiceList

type KuberSourceFactory struct {
	KuberMaster string
	fetchSourceSync
}

type fetchSourceSync struct {
	podStopCh     chan struct{}
	serviceStopCh chan struct{}
	nodeStopCh    chan struct{}
}

//
func NewSourceFactory(kuberMaster string) *KuberSourceFactory {
	return &KuberSourceFactory{
		KuberMaster: kuberMaster,
		fetchSourceSync: fetchSourceSync{
			podStopCh:     make(chan struct{}),
			serviceStopCh: make(chan struct{}),
			nodeStopCh:    make(chan struct{}),
		},
	}
}

//
func (k *KuberSourceFactory) prepare(resource interface{}, path string) {
	defer func() { recover() }()
	var url string
	if strings.Contains(path, "http") {
		url = path
	} else {
		url = k.KuberMaster + path
	}
	if body, err2 := ioutil.ReadAll(func() *http.Response {
		resp, err1 := http.Get(url)
		if err1 != nil {
			util.Logger.SetFatal(err1)
		}
		return resp
	}().Body); err2 != nil {
		util.Logger.SetFatal(err2)
	} else {
		err3 := json.Unmarshal(body, &resource)
		if err3 != nil {
			util.Logger.SetFatal(err2)
		}
	}
	util.Logger.SetDebug("Get done with url:", url)
}

//
func (k *KuberSourceFactory) fetchKuberPod(stopCh chan<- struct{}) {
	k.prepare(&PodList, "/api/v1/pods")
	storage.SourceStorage.Put("pod", PodList)
	stopCh <- struct{}{}
}

//
func (k *KuberSourceFactory) fetchKuberNode(stopCh chan<- struct{}) {
	syncStop := sync.WaitGroup{}
	k.prepare(&NodeList, "/api/v1/nodes")
	syncStop.Add(len(NodeList.Items) * 2)
	for _, v := range NodeList.Items {
		go func(v modelK8s.Node) {
			k.fetchKuberNodePS(&syncStop, v.Status.Addresses[1].Address)

		}(v)
		go func(v modelK8s.Node) {
			k.fetchKuberNodeFS(&syncStop, v.Status.Addresses[1].Address)

		}(v)
	}
	storage.SourceStorage.Put("node", NodeList)
	syncStop.Wait()
	stopCh <- struct{}{}
	return
}

//
func (k *KuberSourceFactory) fetchKuberNodePS(stopSy *sync.WaitGroup, address string) (cpu float32, mem float32) {
	defer func() {
		recover()
		stopSy.Done()
	}()
	var y []v2.ProcessInfo
	if address == "127.0.0.1" {
		panic("node IP could not set 127.0.0.1")
	}
	k.prepare(&y, "http://"+address+":4194/api/v2.0/ps/")
	var c, m float32
	for _, v := range y {
		c = c + v.PercentCpu
		m = m + v.PercentMemory
	}
	cpu = c
	mem = m
	storage.SourceStorage.Put("cpu"+" "+address, cpu)
	storage.SourceStorage.Put("mem"+" "+address, mem)
	return
}

//
func (k *KuberSourceFactory) fetchKuberNodeFS(stopSy *sync.WaitGroup, address string) (int64, int64) {
	defer func() {
		recover()
		stopSy.Done()
	}()
	if address == "127.0.0.1" {
		panic("node IP could not set 127.0.0.1")
	}
	var fs []v2.MachineFsStats
	k.prepare(&fs, "http://"+address+":4194/api/v2.0/storage")
	var outCapacity uint64
	var outUse uint64
	for _, v := range fs {
		outCapacity = *v.Capacity + outCapacity
		outUse = *v.Usage + outUse
	}
	storage.SourceStorage.Put("totalFs"+" "+address, int64(outCapacity))
	storage.SourceStorage.Put("useFS"+" "+address, int64(outUse))
	return int64(outCapacity), int64(outUse)
}

//
func (k *KuberSourceFactory) FetchResource() {
	go k.fetchKuberPod(k.podStopCh)
	go k.fetchKuberService(k.serviceStopCh)
	go k.fetchKuberNode(k.nodeStopCh)
	k.fetchDone()
}

//
func (k *KuberSourceFactory) fetchDone() {
	<-k.podStopCh
	<-k.serviceStopCh
	<-k.nodeStopCh
}

//
func (k *KuberSourceFactory) fetchKuberService(stopCh chan<- struct{}) {
	k.prepare(&ServiceList, "/api/v1/services")
	storage.SourceStorage.Put("service", ServiceList)
	stopCh <- struct{}{}
}
