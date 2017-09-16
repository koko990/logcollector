package cache

import (
	"sync"
	"time"
)

type GenericCache interface {
	Get(key string) interface{}
	GetMulti(keys []string) []interface{}
	Put(key string, val interface{}) error
	Delete(key string) error
	IsExist(key string) bool
	MatchQueryKey(keyMatch string) []interface{}
	ClearAll() error
}

type SourceItem struct {
	val         interface{}
	createdTime time.Time
}

type SourceCache struct {
	key interface{}
	sync.RWMutex
	dur   time.Duration
	items map[string]*SourceItem
}
