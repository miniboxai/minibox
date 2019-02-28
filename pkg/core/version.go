package core

import (
	"errors"
	"fmt"
)

var versions = make(map[string]interface{})
var ErrNonConfigStruct = errors.New("invalid registred config struct")

func Registry(ver string, stuc interface{}) {
	if _, ok := versions[ver]; ok {
		panic(fmt.Errorf("exists registry the %s version ", ver))
	}

	versions[ver] = stuc
}

func getConfigStruct(ver string) (interface{}, error) {
	if cfgStruct, ok := versions[ver]; ok {
		return cfgStruct, nil
	}
	return nil, ErrNonConfigStruct
}
