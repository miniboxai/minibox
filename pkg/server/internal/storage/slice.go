package storage

import (
	"errors"
	"reflect"
)

type Slice struct {
	Val      reflect.Value
	EleTyp   reflect.Type
	convFrom *convertTable
	target   *reflect.Value
}

func fromSlice(vals interface{}) *Slice {
	var (
		val    = reflect.ValueOf(vals)
		eleTyp = guessMember(vals)
	)

	from, ok := findType(eleTyp)
	if !ok {
		panic(errors.New("invalid from convert table"))
	}

	return &Slice{
		Val:      val,
		EleTyp:   eleTyp,
		convFrom: from,
		// convTo:
	}
}

func (s *Slice) Type() reflect.Type {
	return s.EleTyp
}

func (s *Slice) Target() reflect.Value {
	if s.target == nil {
		v := s.convFrom.MakeSlice(0, 0)
		s.target = &v
	}

	return *s.target
}

func (s *Slice) TargetType() reflect.Type {
	return s.convFrom.Target()
}

func (s *Slice) Save() bool {
	if s.target == nil {
		return false
	}
	st := reflect.Indirect(*s.target)
	v := s.Val
	ss := reflect.SliceOf(s.EleTyp)
	v.Elem().Set(reflect.MakeSlice(ss, st.Len(), st.Len()))

	for i := 0; i < st.Len(); i++ {
		v := v.Elem().Index(i)

		if v.Kind() == reflect.Ptr && v.IsNil() {
			v.Set(reflect.New(s.EleTyp))
		}
		v = reflect.Indirect(v)
		e := st.Index(i)
		s.convFrom.ConvertVal(v, e.Addr())
	}

	return true
}
