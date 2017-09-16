package control

import (
	"errors"
	"net/http"
	"reflect"
	"strings"
	"fmt"
)

var routerMap map[string]Router
var funcRepo map[string]repoType

type repoType struct {
	funcName reflect.Value
	method   string
}

func init() {
	routerMap = make(map[string]Router)
	routerMap["getStatus"] = Router{Path: "/{status}", HandlerFunc: getStatus, Method: "POST"}
	routerMap["getStatusIndex"] = Router{Path: "/", HandlerFunc: getStatusIndex, Method: "GET"}
}

type Api struct {
}

type Router struct {
	Path        string
	HandlerFunc http.HandlerFunc
	Method      string
}

func (Api)setRouterCfg(rp string) error {
	var (
		methodName string
		funcString string
	)
	stringTemp := strings.Split(rp, "_")
	if len(stringTemp) != 2 {
		return errors.New("routerMap is wrong")
	}
	methodName = stringTemp[0]
	funcString = stringTemp[1]
	aa := reflect.ValueOf(Api{}).MethodByName(funcString)
	if aa.IsValid() {
		return errors.New("router fun is wrong")
	}
	funcRepo[funcString] = repoType{
		funcName: aa,
		method:   methodName,
	}
	return nil

}
func (Api)fetchFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;   charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "post the app status")
	v := make([]reflect.Value, 0)
	for _,val:=range funcRepo{
		val.funcName.Call(v)
	}
}
