package storage

import (
	"errors"
	"reflect"
)

type Value struct {
	Val      reflect.Value
	Typ      reflect.Type
	convFrom *convertTable
	convTo   *convertTable
	target   *reflect.Value
}

func from(val interface{}) *Value {
	var (
		v = reflect.ValueOf(val)
		t = reflect.Indirect(v).Type()
	)

	from, ok := findType(t) // types
	if !ok {
		panic(errors.New("invalid from convert table"))
	}

	to, ok := findType(from.Target())
	if !ok {
		panic(errors.New("invalid to convert table"))
	}

	return &Value{
		Val:      v,
		Typ:      t,
		convFrom: from,
		convTo:   to,
	}
}

func (v *Value) Target() reflect.Value {
	if v.target == nil {
		nv := v.convFrom.NewTarget()
		v.target = &nv // models
		v.convTo.ConvertVal(v.target.Elem(), v.Val)
	}

	return *v.target
}

func (v *Value) TargetType() reflect.Type {
	return v.convFrom.Target()
}

func (v *Value) Type() reflect.Type {
	return v.Typ
}

func (v *Value) Save() bool {
	if v.target == nil {
		return false
	}

	v.convFrom.ConvertVal(v.Val.Elem(), *v.target)
	return true
}
