package dao

import (
	"github.com/koko990/logcollector/models/types"

	"context"
	"github.com/koko990/logcollector/models/util/pool"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

const (
	pushDBTime = time.Second * 5
)

//
func PushMultiModelToDB(ctx context.Context, cancel context.CancelFunc) {
	defer func() {
		recover()
		cancel()
	}()
	for range time.Tick(pushDBTime) {
		checkModel()
		var serviceValItems *[]types.ServiceDashboard
		var nodeValItems *[]types.DashboardNode
		serviceValItems = new([]types.ServiceDashboard)
		nodeValItems = new([]types.DashboardNode)
		done1Ch := make(chan struct{})
		done2Ch := make(chan struct{})
		go fetchServiceFromPool(serviceValItems, done1Ch)
		go fetchNodeFromPool(nodeValItems, done2Ch)
		<-done1Ch
		<-done2Ch
		if n := len(*serviceValItems); n != 0 {
			insertMultiModel(n, serviceValItems)

		}
		if n := len(*nodeValItems); n != 0 {
			insertMultiModel(n, nodeValItems)
		}

	}

}

//
func fetchServiceFromPool(serviceValItems *[]types.ServiceDashboard, done1Ch chan<- struct{}) {
	for pool.ServiceDAOBuffer.Back() != nil {
		val := pool.ServiceDAOBuffer.Back().Value
		serviceVal := val.(types.ServiceDashboard)
		*serviceValItems = append(*serviceValItems, serviceVal)
		pool.ServiceDAOBuffer.Remove(pool.ServiceDAOBuffer.Back())
	}
	done1Ch <- struct{}{}
}

//
func fetchNodeFromPool(nodeValItems *[]types.DashboardNode, done2Ch chan<- struct{}) {
	for pool.NodeDAOBuffer.Back() != nil {
		val := pool.NodeDAOBuffer.Back().Value
		nodeVal := val.(types.DashboardNode)
		*nodeValItems = append(*nodeValItems, nodeVal)
		pool.NodeDAOBuffer.Remove(pool.NodeDAOBuffer.Back())

	}
	done2Ch <- struct{}{}
}

//
func insertMultiModel(len int, model interface{}) (int64, error) {
	o := orm.NewOrm()
	id, err := o.InsertMulti(len, model)
	return id, err
}
