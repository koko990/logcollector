package merge

import (
	"github.com/koko990/logcollector/models/util/pool"
	"github.com/koko990/logcollector/service/storage"

	"fmt"
	"github.com/koko990/logcollector/models/types"
	"reflect"
	"strconv"
	"sync"
)

func RunNode(timestampArg string, nodeAddress string, m *sync.Mutex) {
	if checkNodeCache(timestampArg, nodeAddress) == true {
		var (
			nodeTemp  *types.DashboardNode
			cpuCh     chan struct{} = make(chan struct{})
			memoryCh  chan struct{} = make(chan struct{})
			totalFsCh chan struct{} = make(chan struct{})
			usageFsCh chan struct{} = make(chan struct{})
		)
		nodeTemp = &types.DashboardNode{
			NodeName: nodeAddress,
			Timestamp: func(timestampArg string) int {
				t, _ := strconv.Atoi(timestampArg)
				return t
			}(timestampArg),
			InternalIp: nodeAddress,
		}
		go setCpu(timestampArg, nodeAddress, m, nodeTemp, cpuCh)
		go setMemoryUsages(timestampArg, nodeAddress, m, nodeTemp, memoryCh)
		go setTotalFs(timestampArg, nodeAddress, m, nodeTemp, totalFsCh)
		go setUsageFs(timestampArg, nodeAddress, m, nodeTemp, usageFsCh)
		finalTask(timestampArg, nodeTemp, cpuCh, memoryCh, totalFsCh, usageFsCh)
	} else {
	}
}

//
func finalTask(timestampArg string, nodeTemp *types.DashboardNode, cpuCh <-chan struct{},
	memoryCh <-chan struct{}, totalFsCh <-chan struct{}, usageFsCh <-chan struct{}) {
	<-cpuCh
	<-memoryCh
	<-totalFsCh
	<-usageFsCh
	delete(pool.MergeNodeCache, pool.IndexNode{
		Timestamp: timestampArg,
		NodeKey:   int(storage.Node)})
	pool.NodeDAOBuffer.PushBack(*nodeTemp)
}

//
func checkNodeCache(timestampArg string, nodeAddress string) bool {
	if pool.MergeNodeCache[pool.IndexNode{
		Timestamp:   timestampArg,
		NodeKey:     int(storage.Cpu),
		NodeAddress: nodeAddress}] != nil &&
		pool.MergeNodeCache[pool.IndexNode{
			Timestamp:   timestampArg,
			NodeKey:     int(storage.MemoryUsage),
			NodeAddress: nodeAddress}] != nil &&
		pool.MergeNodeCache[pool.IndexNode{
			Timestamp:   timestampArg,
			NodeKey:     int(storage.TotalFs),
			NodeAddress: nodeAddress}] != nil &&
		pool.MergeNodeCache[pool.IndexNode{
			Timestamp:   timestampArg,
			NodeKey:     int(storage.UsageFs),
			NodeAddress: nodeAddress}] != nil {
		return true
	}
	return false
}

//
func setCpu(timestampArg string, nodeAddress string, m *sync.Mutex,
	nodeTemp *types.DashboardNode, done chan<- struct{}) {
	m.Lock()
	indexCpu := pool.IndexNode{
		Timestamp:   timestampArg,
		NodeKey:     int(storage.Cpu),
		NodeAddress: nodeAddress}
	valCpu := pool.MergeNodeCache[indexCpu]
	temp, ok := valCpu.(float32)
	nodeTemp.CpuUsage = temp
	delete(pool.MergeNodeCache, indexCpu)
	m.Unlock()
	if ok == false {
		panic(fmt.Sprint(ok, reflect.TypeOf(temp)))
		done <- struct{}{}
		return
	} else {
		done <- struct{}{}
		return
	}
}

//
func setMemoryUsages(timestampArg string, nodeAddress string,
	m *sync.Mutex, nodeTemp *types.DashboardNode, done chan<- struct{}) {
	m.Lock()
	indexMem := pool.IndexNode{
		Timestamp:   timestampArg,
		NodeKey:     int(storage.MemoryUsage),
		NodeAddress: nodeAddress}
	valMem := pool.MergeNodeCache[indexMem]
	temp, ok := valMem.(float32)
	nodeTemp.MemUsage = temp
	delete(pool.MergeNodeCache, indexMem)
	m.Unlock()
	if ok == false {
		panic(fmt.Sprint(ok, reflect.TypeOf(temp)))
		done <- struct{}{}
		return
	} else {
		done <- struct{}{}
		return
	}
}

//
func setTotalFs(timestampArg string, nodeAddress string, m *sync.Mutex,
	nodeTemp *types.DashboardNode, done chan<- struct{}) {
	m.Lock()
	indexTotalFs := pool.IndexNode{
		Timestamp:   timestampArg,
		NodeKey:     int(storage.TotalFs),
		NodeAddress: nodeAddress}
	valTotal := pool.MergeNodeCache[indexTotalFs]
	temp, ok := valTotal.(int64)
	nodeTemp.StorageTotal = temp
	delete(pool.MergeNodeCache, indexTotalFs)
	m.Unlock()
	if ok == false {
		panic(fmt.Sprint(ok, reflect.TypeOf(valTotal)))
		done <- struct{}{}
		return
	} else {
		done <- struct{}{}
		return
	}
}

//
func setUsageFs(timestampArg string, nodeAddress string, m *sync.Mutex,
	nodeTemp *types.DashboardNode, done chan<- struct{}) {
	m.Lock()
	indexUsageFs := pool.IndexNode{
		Timestamp:   timestampArg,
		NodeKey:     int(storage.UsageFs),
		NodeAddress: nodeAddress}
	valUsageFs := pool.MergeNodeCache[indexUsageFs]
	temp, ok := valUsageFs.(int64)
	nodeTemp.StorageUse = temp
	m.Unlock()
	if ok == false {
		panic(fmt.Sprint(ok, reflect.TypeOf(temp)))
		done <- struct{}{}
		return
	} else {
		done <- struct{}{}
		return
	}

}
