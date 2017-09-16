package merge

import (
	"context"
	"github.com/koko990/logcollector/models/util/pool"
	"github.com/koko990/logcollector/service/storage"
	"sync"
	"time"
)

const (
	defaultMergeTime = time.Millisecond * 100
)

//
func RunMergeModel(ctx context.Context, m *sync.Mutex,cancel context.CancelFunc) {
	for range time.Tick(defaultMergeTime) {
		timestampIndex, add := fetchSourceToPool(m)
		RunService(timestampIndex)
		RunNode(timestampIndex, add, m)
	}
	cancel()

}

//
func fetchSourceToPool(m *sync.Mutex) (string, string) {
	key, indexKey, timestampIndex, addIndex, val, _ := storage.SourceStorage.Get()
	if key == storage.Pod || key == storage.Service {
		m.Lock()
		pool.MergeServiceCache[pool.IndexService{
			Timestamp:  timestampIndex,
			ServiceKey: int(indexKey)}] = val
		m.Unlock()
	} else {
		m.Lock()
		pool.MergeNodeCache[pool.IndexNode{
			Timestamp:   timestampIndex,
			NodeKey:     int(indexKey),
			NodeAddress: addIndex}] = val
		m.Unlock()
	}
	return timestampIndex, addIndex
}
