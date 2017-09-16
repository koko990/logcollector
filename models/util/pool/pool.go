package pool

import "container/list"

type IndexNode struct {
	Timestamp   string
	NodeKey     int
	NodeAddress string
}

type IndexService struct {
	Timestamp  string
	ServiceKey int
}

var MergeNodeCache map[IndexNode]interface{}
var MergeServiceCache map[IndexService]interface{}
var ServiceDAOBuffer list.List
var NodeDAOBuffer list.List

func init() {
	MergeNodeCache = make(map[IndexNode]interface{})
	MergeServiceCache = make(map[IndexService]interface{})
	ServiceDAOBuffer.Init()
	NodeDAOBuffer.Init()

}
