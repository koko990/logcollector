package cache

import (
	"bytes"
	"errors"
	"time"
)

func NewMemoryCache(key string) *SourceCache {
	newCache := SourceCache{key: key, items: make(map[string]*SourceItem)}
	return &newCache
}

func (bc *SourceCache) Get(name string) interface{} {
	bc.RLock()
	defer bc.RUnlock()
	if itm, ok := bc.items[name]; ok {
		return itm.val
	}
	return nil
}

func (bc *SourceCache) GetMulti(names []string) []interface{} {
	var rc []interface{}
	for _, name := range names {
		rc = append(rc, bc.Get(name))
	}
	return rc
}

func (bc *SourceCache) Put(name string, value interface{}) error {
	bc.Lock()
	defer bc.Unlock()
	bc.items[name] = &SourceItem{
		val:         value,
		createdTime: time.Now(),
	}
	return nil
}

func (bc *SourceCache) Delete(name string) error {
	bc.Lock()
	defer bc.Unlock()
	if _, ok := bc.items[name]; !ok {
		return errors.New("key not exist")
	}
	delete(bc.items, name)
	if _, ok := bc.items[name]; ok {
		return errors.New("delete key error")
	}
	return nil
}

func (bc *SourceCache) IsExist(name string) bool {
	bc.RLock()
	defer bc.RUnlock()
	_, ok := bc.items[name]
	return ok
}

func (bc *SourceCache) ClearAll() error {
	bc.Lock()
	defer bc.Unlock()
	bc.items = make(map[string]*SourceItem)
	return nil
}

func (bc *SourceCache) clearItems(keys []string) {
	bc.Lock()
	defer bc.Unlock()
	for _, key := range keys {
		delete(bc.items, key)
	}
}

func (bc *SourceCache) MatchQueryKey(keyMatch string) []interface{} {
	var rc []interface{}
	for k := range bc.items {
		if bytes.Contains([]byte(k), []byte(keyMatch)) {
			rc = append(rc, bc.Get(k))
		}
	}
	return rc
}

func Register(repoName string) GenericCache {
	return NewMemoryCache(repoName)
}
