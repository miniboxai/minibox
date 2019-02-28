package utils

import (
	"reflect"
	"strings"
	"time"

	oerrors "errors"

	"minibox.ai/minibox/pkg/core/errors"
)

type Map = map[interface{}]interface{}
type Any = interface{}
type VMap map[interface{}]interface{}

func MustSub(m Map, name string) interface{} {
	if v, ok := m[name]; !ok {
		panic(&errors.ErrInvalidConfigItem{Name: name})
	} else {
		return v
	}
}

func MustMap(m Map, name string) Map {
	if v, ok := m[name].(Map); !ok {
		panic(&errors.ErrInvalidConfigItem{Name: name})
	} else {
		return v
	}
}

func MustStringSlice(m Any, name string) []string {
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

func MustString(m Any, name string) string {
	if v, ok := m.(string); !ok {
		panic(&errors.ErrInvalidItemType{Name: name, Type: reflect.TypeOf("")})
	} else {
		return v
	}
}

func MustInt(m Any, name string) int {
	if v, ok := m.(int); !ok {
		panic(&errors.ErrInvalidItemType{Name: name, Type: reflect.TypeOf(0)})
	} else {
		return v
	}
}

func MustBool(m Any, name string) bool {
	if v, ok := m.(bool); !ok {
		panic(&errors.ErrInvalidItemType{Name: name, Type: reflect.TypeOf(false)})
	} else {
		return v
	}
}

func (m VMap) Map(name string) VMap {
	if v, ok := m[name].(map[interface{}]interface{}); !ok {
		panic(&errors.ErrInvalidConfigItem{Name: name})
	} else {
		return VMap(v)
	}
}

func (m VMap) String(name string) string {
	if v, ok := m[name].(string); !ok {
		panic(&errors.ErrInvalidItemType{Name: name, Type: reflect.TypeOf("")})
	} else {
		return v
	}
}

func (m VMap) Int(name string) int {
	if v, ok := m[name].(int); !ok {
		panic(&errors.ErrInvalidItemType{Name: name, Type: reflect.TypeOf("")})
	} else {
		return v
	}
}

func (m VMap) Int16(name string) int16 {
	if v, ok := m[name].(int16); !ok {
		panic(&errors.ErrInvalidItemType{Name: name, Type: reflect.TypeOf(int16(0))})
	} else {
		return v
	}
}

func (m VMap) Int32(name string) int32 {
	if v, ok := m[name].(int32); !ok {
		panic(&errors.ErrInvalidItemType{Name: name, Type: reflect.TypeOf(int32(0))})
	} else {
		return v
	}
}

func (m VMap) Int64(name string) int64 {
	if v, ok := m[name].(int64); !ok {
		panic(&errors.ErrInvalidItemType{Name: name, Type: reflect.TypeOf(int64(0))})
	} else {
		return v
	}
}

func (m VMap) Bool(name string) bool {
	if v, ok := m[name].(bool); !ok {
		panic(&errors.ErrInvalidItemType{Name: name, Type: reflect.TypeOf(false)})
	} else {
		return v
	}
}

func (m VMap) Float64(name string) float64 {
	if v, ok := m[name].(float64); !ok {
		panic(&errors.ErrInvalidItemType{Name: name, Type: reflect.TypeOf(0.0)})
	} else {
		return v
	}
}

func (m VMap) Time(name string) time.Time {
	if v, ok := m[name].(string); !ok {
		panic(&errors.ErrInvalidItemType{Name: name, Type: reflect.TypeOf(time.Time{})})
	} else {
		if t, err := time.Parse(time.RFC3339, v); err != nil {
			panic(err)
		} else {
			return t
		}
	}
}
func (m VMap) StringSlice(name string) []string {
	if arr, ok := m[name].([]interface{}); !ok {
		panic(&errors.ErrInvalidItemType{Name: name, Type: reflect.TypeOf([]string{})})
	} else {
		var ret = make([]string, 0, len(arr))
		for _, v := range arr {
			if v, ok := v.(string); ok {
				ret = append(ret, v)
			}
		}
		return ret
	}
}

func (m VMap) Sub(pat, sep string) VMap {
	ss := strings.Split(pat, sep)
	var v = m
	for _, s := range ss {
		v = v.Map(s)
	}
	return v
}

func (m VMap) Range(handle func(key, val Any) error) error {
	for key, val := range m {
		if err := handle(key, val); err != nil {
			return err
		}
	}
	return nil
}

func (m VMap) RangeSKey(handle func(key string, val Any) error) error {
	for k, val := range m {
		if key, ok := k.(string); ok {
			if err := handle(key, val); err != nil {
				return err
			}
		} else {
			return &errors.ErrKeyIsNotString{k}
		}
	}
	return nil
}

func Any2Map(a Any) VMap {
	if m, ok := a.(map[interface{}]interface{}); !ok {
		panic(oerrors.New("is invalid map[interface{}]interface{} Type"))
	} else {
		return VMap(m)
	}
}
