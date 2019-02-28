package storage

import (
	"reflect"
	"sync"
)

type converts struct {
	sync.Map
}

type convertTable struct {
	target   reflect.Type
	convFunc reflect.Value
}

var globalConverts = &converts{}

func convFunc(convFun interface{}) reflect.Value {
	fun := reflect.ValueOf(convFun)
	if fun.Kind() != reflect.Func {
		panic("convFunc must Func")
	}

	return fun
}

func typeFunc(a interface{}) reflect.Type {
	v := reflect.ValueOf(a)
	t := reflect.Indirect(v).Type()
	return t
}

func find(a interface{}) (*convertTable, bool) {
	t := typeFunc(a)
	if v, ok := globalConverts.Load(t.String()); ok {
		v, ok := v.(*convertTable)
		return v, ok
	}

	return nil, false
}

func findType(t reflect.Type) (*convertTable, bool) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if v, ok := globalConverts.Load(t.String()); ok {
		v, ok := v.(*convertTable)
		return v, ok
	}
	return nil, false
}

func RegisterConvert(a, b interface{}, conv, reverse interface{}) {
	t1 := typeFunc(a)
	t2 := typeFunc(b)

	globalConverts.Store(t1.String(), &convertTable{t2, convFunc(reverse)})
	globalConverts.Store(t2.String(), &convertTable{t1, convFunc(conv)})
}

func (c *convertTable) Target() reflect.Type {
	return c.target
}

func (c *convertTable) NewTarget() reflect.Value {
	return reflect.New(c.target)
}

func (c *convertTable) MakeSlice(l, ca int) reflect.Value {
	t := reflect.SliceOf(c.target)
	return reflect.New(t)
	// return reflect.MakeSlice(t, l, ca)
}

func (c *convertTable) Convert(dest interface{}, src interface{}) {
	v := reflect.ValueOf(src)
	res := c.convFunc.Call([]reflect.Value{v})
	vv := reflect.ValueOf(dest)
	vv.Set(res[0])
}

func (c *convertTable) ConvertVal(dest, src reflect.Value) {
	var (
		out = c.convFunc.Call([]reflect.Value{src})
		val = out[0]
	)
	dest.Set(val.Elem())
}

func (c *convertTable) ConvertSlice(dest, src interface{}) {
	// list, ok := src.([]interface{})
	// if !ok {
	// 	panic("src is not slice")
	// }

	// for i, m := range list {
	// 	c.convert()
	// }
}
