package dao

import (
	"fmt"
	"github.com/koko990/logcollector/models/types"
	"time"

	"github.com/koko990/logcollector/util"

	"github.com/astaxie/beego/orm"
)

var currentTime string
var prepareStatus bool

func pushInit() {
	defer func() {
		recover()
	}()
	orm.RegisterModel(new(types.ServiceDashboard), new(types.DashboardNode))
	fmt.Println("Initializing DB registration.")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", "root:90660652@tcp(localhost:3306)/board?charset=utf8")
	orm.RunSyncdb("default", false, true)
	if err != nil {
		util.Logger.SetError("Error occurred on registering DB: %+v\n", err)
	}
}
func checkModel() {
	tmpTime := time.Now().Format("2006 01 02")
	if currentTime != tmpTime || prepareStatus == false {
		orm.ResetModelCache()
		pushInit()
		prepareStatus = true
		currentTime = tmpTime
	}
}
