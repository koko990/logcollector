package storage

import (
	"container/list"
	"github.com/koko990/logcollector/util"
	"strings"
	"unicode"
	"sync"
)

var TimestampIndex string
var SourceStorage storeList
var valStore map[string]interface{}

type Key int

const (
	_       Key = iota
	Pod
	Node
	Service
	Cpu
	MemoryUsage
	TotalFs
	UsageFs
)

type storeList struct {
	e *list.List
	m *sync.Mutex
}

//
func (s *storeList) Put(key string, val interface{}) {
	s.m.Lock()
	defer func(*sync.Mutex) {
		s.m.Unlock()
	}(s.m)
	key = key + " " + TimestampIndex
	s.e.PushBack(key)
	valStore[key] = val
}

//
func (s *storeList) Get() (key Key, index Key, timestamp string, add string, val interface{}, err error) {
	var mapKey string
	if s.e.Back() != nil {
		mapKey = (s.e.Back().Value).(string)
		timestamp = strings.FieldsFunc(mapKey, unicode.IsSpace)[len(strings.FieldsFunc(mapKey, unicode.IsSpace))-1]
		util.Logger.SetDebug("mapKey is ", mapKey, s.e.Len())
	} else {
		return 0, 0, timestamp, "", nil, nil
	}
	if strings.Contains(mapKey, "pod") {
		return Pod, Pod, timestamp, "", s.deleteSpecificListMap(mapKey), nil
	}
	if strings.Contains(mapKey, "node") {
		return Node, Node, timestamp, "", s.deleteSpecificListMap(mapKey), nil
	}
	if strings.Contains(mapKey, "cpu") {
		return Node, Cpu, timestamp, strings.FieldsFunc(mapKey, unicode.IsSpace)[1], s.deleteSpecificListMap(mapKey), nil
	}
	if strings.Contains(mapKey, "mem") {
		return Node, MemoryUsage, timestamp, strings.FieldsFunc(mapKey, unicode.IsSpace)[1], s.deleteSpecificListMap(mapKey), nil
	}
	if strings.Contains(mapKey, "totalFs") {
		return Node, TotalFs, timestamp, strings.FieldsFunc(mapKey, unicode.IsSpace)[1], s.deleteSpecificListMap(mapKey), nil
	}
	if strings.Contains(mapKey, "useFS") {
		return Node, UsageFs, timestamp, strings.FieldsFunc(mapKey, unicode.IsSpace)[1], s.deleteSpecificListMap(mapKey), nil
	}
	if strings.Contains(mapKey, "service") {
		return Service, Service, timestamp, "", s.deleteSpecificListMap(mapKey), nil
	}

	return
}

//
func (s *storeList) deleteSpecificListMap(mapKey string) interface{} {
	s.m.Lock()
	s.e.Remove(s.e.Back())
	val := valStore[mapKey]
	delete(valStore, mapKey)
	s.m.Unlock()
	return val
}
func (s *storeList) InitStore(mutex *sync.Mutex) {
	valStore = make(map[string]interface{})
	SourceStorage.e = list.New()
	SourceStorage.m = mutex
}

//
