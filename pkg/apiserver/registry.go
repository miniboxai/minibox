package apiserver

import (
	"errors"
	"sort"
)

var ErrAlwaysRegisterdApiServer = errors.New("always registered ApiServer")

type ApiInitialFunc func(*ApiServer) error

type apiIntial struct {
	idx    int
	handle ApiInitialFunc
}

type initails []apiIntial

func (a initails) Len() int           { return len(a) }
func (a initails) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a initails) Less(i, j int) bool { return a[i].idx < a[j].idx }

var registryApiServer *ApiServer

// var initializeFuncs = make([]ApiInitialFunc, 0)
var initializeFuncs = make([]apiIntial, 0)

func Initializer(idx int, fn ApiInitialFunc) {
	initializeFuncs = append(initializeFuncs, apiIntial{idx, fn})
}

func Registry(svr *ApiServer) error {
	if registryApiServer != nil {
		return ErrAlwaysRegisterdApiServer
	}

	registryApiServer = svr
	return nil
}

func EachInitializer(svr *ApiServer) {
	initializeFuncs := initails(initializeFuncs)
	sort.Sort(initializeFuncs)
	for _, aInt := range initializeFuncs {
		aInt.handle(svr)
	}
}
