package core

import (
	"reflect"

	"minibox.ai/minibox/pkg/core/errors"
)

type Map = map[interface{}]interface{}
type Any = interface{}

func mustSub(m Map, name string) interface{} {
	if v, ok := m[name]; !ok {
		panic(&errors.ErrInvalidConfigItem{Name: name})
	} else {
		return v
	}
}

func mustMap(m Map, name string) Map {
	if v, ok := m[name].(Map); !ok {
		panic(&errors.ErrInvalidConfigItem{Name: name})
	} else {
		return v
	}
}

func mustStringSlice(m Any, name string) []string {
	if vs, ok := m.([]interface{}); !ok {
		panic(&errors.ErrInvalidItemType{Name: name, Type: reflect.TypeOf([]string{})})
	} else {
		var ss = make([]string, 0, len(vs))
		for _, v := range vs {
			if s, ok := v.(string); !ok {
				panic(&errors.ErrInvalidItemType{Name: name, Type: reflect.TypeOf([]string{})})
			} else {
				ss = append(ss, s)
			}

		}
		return ss
	}
}

func mustString(m Any, name string) string {
	if v, ok := m.(string); !ok {
		panic(&errors.ErrInvalidItemType{Name: name, Type: reflect.TypeOf("")})
	} else {
		return v
	}
}

func mustInt(m Any, name string) int {
	if v, ok := m.(int); !ok {
		panic(&errors.ErrInvalidItemType{Name: name, Type: reflect.TypeOf(0)})
	} else {
		return v
	}
}

func mustBool(m Any, name string) bool {
	if v, ok := m.(bool); !ok {
		panic(&errors.ErrInvalidItemType{Name: name, Type: reflect.TypeOf(false)})
	} else {
		return v
	}
}
