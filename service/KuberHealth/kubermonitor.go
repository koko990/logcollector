package KuberHealth

import (
	"fmt"
	"github.com/koko990/logcollector/util"
	"io/ioutil"
	"net/http"
	"time"
)

func KuberMonitor(kuberMaster string, status *bool) {
	defer die(kuberMaster)
	for range time.Tick(time.Second) {
		url := fmt.Sprintf("%s/version", kuberMaster)
		cl := &http.Client{Timeout: time.Millisecond * 2000}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}
		resp, _ := cl.Do(req)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		if len(body) > 1 {
			util.Logger.SetInfo("kubernetes version is run")
			*status = true
		} else {
			util.Logger.SetError("kuber is down", kuberMaster)
			panic("kuber is down")
		}

	}

}

func die(kuberMaster string) {
	recover()
	util.Logger.SetError("kuber is die,kuber url: ", kuberMaster)
}
