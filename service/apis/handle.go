package control

import (
	"errors"
	"fmt"
	"github.com/koko990/logcollector/util"
	"net/http"

	"github.com/gorilla/mux"
)

func CollectRouters() (router *mux.Router, err error) {
	router = mux.NewRouter().StrictSlash(true).PathPrefix("/status").Subrouter() /*StrictSlash: /path/ to /path */
	if router == nil {
		return nil, err
	}
	for key, v := range routerMap {
		var handler http.Handler
		handler = util.Logger.HttpLog(v.HandlerFunc, key)
		router.Methods(v.Method).Path(v.Path).Handler(handler)
		if handler == nil {
			err = errors.New("func is wrong" + key)
			return nil, err
		}
		if err != nil {
			return nil, err
		}
	}
	return router, nil
}

func responseHead(w http.ResponseWriter, r *http.Request, bodyString string) {
	w.Header().Set("Content-Type", "application/json;   charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, bodyString)
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars["status"]
	util.Logger.SetInfo("collect")
	switch status {
	case "false":
		//app.TurnStatus(false)
		//*app.Switch = false
		responseHead(w, r, "post turn off")
	case "0":
		//*app.Switch = false
		//app.TurnStatus(false)
		responseHead(w, r, "post turn off")
	default:
		//*app.Switch = true
		//app.TurnStatus(true)
		responseHead(w, r, "post turn run")

	}

}

func getStatusIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;   charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "post the app status")
	util.Logger.SetInfo("collect")
}
