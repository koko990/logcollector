package app

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/koko990/logcollector/cmd/app/options"
	"github.com/koko990/logcollector/models"
	"github.com/koko990/logcollector/service"
	"github.com/koko990/logcollector/service/apis"
	"github.com/koko990/logcollector/service/storage"
	"github.com/koko990/logcollector/util"
)

const storageRebootTime = time.Millisecond * 50

func CollectorDaemon() {
	var err error
	han := func() http.Handler {
		s, _ := control.CollectRouters()
		return s
	}()
	go func() { err = http.ListenAndServe(":8087", han) }()
	svr := service.CollectorServerFactory(setCollectorServer())
	var signalCh chan struct{}
	signalCh = make(chan struct{})
	m := new(sync.Mutex)
	storage.SourceStorage.InitStore(m)
	var storeJob = models.StoreJob{}
	go signal(signalCh)
	go svr.RunServer(signalCh)
	go func() {
		for range time.Tick(storageRebootTime) {
			storeJob.Run(m)
		}

	}()
	if err != nil {
		util.Logger.SetFatal(err)
		return
	}

}

func signal(signCh chan struct{}) {
	for range signCh {
		util.Logger.SetInfo("logcollector  status is correct")
	}
}
func setCollectorServer() string {
	kuberIp := options.RunFlag.ServerKubeIp
	kuberPort := options.RunFlag.ServerKubePort
	return fmt.Sprintf("http://%s:%s", kuberIp, kuberPort)

}
