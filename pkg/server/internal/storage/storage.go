package storage

import (
	"errors"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
	"minibox.ai/minibox/pkg/server/internal/option"
	"minibox.ai/minibox/pkg/utils"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Load(value interface{}, opts ...option.OptionFunc) error {
	var (
		val = reflect.ValueOf(value)
		t   = reflect.Indirect(val).Type()
		opt *option.Option
	)
	opt = loadTypeOptions(t, opts)
	from, ok := findType(t) // types

	if !ok {
		return errors.New("invalid from convert table")
	}

	to, ok := findType(from.Target())
	if !ok {
		return errors.New("invalid to convert table")
	}

	v := from.NewTarget() // models
	vv := v.Interface()

	to.ConvertVal(v.Elem(), val)

	db := s.db
	if len(opt.SubLoads) > 0 {
		for _, sub := range opt.SubLoads {
			db = db.Preload(sub.Name, sub.Args...)
		}
	}

	if opt.Undeleted {
		db = db.Unscoped()
	}

	// db = db.Related(&models.User{}, "users")
	if err := db.Find(vv).Error; err != nil {
		return err
	}

	from.ConvertVal(val.Elem(), v)
	return nil
}

func (s *Storage) LoadOrStore(value interface{}) (interface{}, error) {
	return nil, nil
}

func (s *Storage) updateAttrs(uds []option.UpdateField, vv reflect.Value) {
	v := reflect.Indirect(vv)
	if v.Kind() != reflect.Struct {
		panic(ErrMustStructValue)
	}

	for _, ud := range uds {
		fv := v.FieldByName(ud.Name)
		fv.Set(ud.Value)
	}
}

func (s *Storage) Store(value interface{}, opts ...option.OptionFunc) error {
	var (
		val = reflect.ValueOf(value)
		t   = reflect.Indirect(val).Type()
		opt *option.Option
	)

	opt = loadTypeOptions(t, opts)
	from, ok := findType(t) // types

	if !ok {
		return errors.New("invalid from convert table")
	}

	to, ok := findType(from.Target())
	if !ok {
		return errors.New("invalid to convert table")
	}

	v := from.NewTarget() // models
	vv := v.Interface()
	fieldKey, ok := s.primaryKey(from.Target())
	if !ok {
		return errors.New("invalid primary key")
	}

	to.ConvertVal(v.Elem(), val)

	if len(opt.Updates) > 0 {
		s.updateAttrs(opt.Updates, v)
	}

	prkey := v.Elem().FieldByName(fieldKey.Name)

	db := s.db

	if opt.Undeleted {
		db = db.Unscoped()
	}

	// 空值, 意味着需生成一个 ID
	if utils.IsZero(prkey) {
		// 如果是字符串类型, 判断是 auto Tag 是什么类型
		// 自动生成 UUID
		if fieldKey.Type.Kind() == reflect.String {
			if fieldKey.Tag.Get("auto") == "uuid" {
				prkey.SetString(utils.UUID())
			}
		}
		if err := db.Create(vv).Error; err != nil {
			return err
		}
	} else if opt.AssignAttrs != nil {
		if err := db.Model(vv).Updates(opt.AssignAttrs).Error; err != nil {
			return err
		}
	} else {
		if err := db.Model(vv).Updates(vv).Error; err != nil {
			return err
		}
	}

	from.ConvertVal(val.Elem(), v)

	return nil
}

// 返回主键结构字段, 如果没有设置 gorm:"primary_key" , 那么就返回
// 近似 ID 的字段
func (s *Storage) primaryKey(t reflect.Type) (reflect.StructField, bool) {
	if t.Kind() != reflect.Struct {
		panic(errors.New("must struct type."))
	}

	var id *reflect.StructField

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			field, ok := s.primaryKey(field.Type)
			if ok {
				return field, ok
			}
		}

		gorm := field.Tag.Get("gorm")
		if strings.Contains(gorm, "primary_key") {
			return field, true
		}

		if strings.ToLower(field.Name) == "id" {
			id = &field
		}
	}

	if id != nil {
		return *id, true
	}

	return reflect.StructField{}, false
}

func (s *Storage) Delete(value interface{}, opts ...option.OptionFunc) error {
	var (
		val = reflect.ValueOf(value)
		t   = reflect.Indirect(val).Type()
		opt *option.Option
	)

	opt = loadTypeOptions(t, opts)
	from, ok := findType(t) // types

	if !ok {
		return errors.New("invalid from convert table")
	}

	to, ok := findType(from.Target())
	if !ok {
		return errors.New("invalid to convert table")
	}

	v := from.NewTarget() // models
	vv := v.Interface()
	fieldKey, ok := s.primaryKey(from.Target())
	if !ok {
		return errors.New("invalid primary key")
	}

	to.ConvertVal(v.Elem(), val)

	prkey := v.Elem().FieldByName(fieldKey.Name)
	// 空值, 抛出无效数据
	if utils.IsZero(prkey) {
		return ErrMustPrimarykeyValue
	}

	db := s.db
	if opt.Force {
		db = db.Unscoped()
	}

	if err := db.Delete(vv).Error; err != nil {
		return err
	}

	from.ConvertVal(val.Elem(), v)

	return nil
}

func (s *Storage) List(list interface{}, opts ...option.OptionFunc) error {
	var (
		t   = guessMember(list)
		sl  interface{}
		st  reflect.Value
		opt *option.Option
	)

	opt = loadOptions(list, opts)
	ct, ok := findType(t)
	// log.Printf()
	if !ok {
		sl = list
	} else {
		st = ct.MakeSlice(0, 0)
		sl = st.Interface()
	}

	db := s.db

	// options apply
	if opt.Limit > 0 {
		db = db.Limit(opt.Limit)
	}

	if opt.Offset > 0 {
		db = db.Offset(opt.Offset)
	}

	if len(opt.SubLoads) > 0 {
		for _, sub := range opt.SubLoads {
			db = db.Preload(sub.Name, sub.Args...)
		}
	}

	if opt.Filters != nil {
		db = opt.Filters.Apply(db)
	}

	if opt.Undeleted {
		db = db.Unscoped()
	}

	// exec sql execute
	if err := db.Find(sl).Error; err != nil {
		return err
	}

	// convertSlice(list, sl)
	st = reflect.Indirect(st)
	if st.Kind() != reflect.Slice {
		panic("src is not slice")
	}

	v := reflect.ValueOf(list)
	ss := reflect.SliceOf(t)
	v.Elem().Set(reflect.MakeSlice(ss, st.Len(), st.Len()))

	for i := 0; i < st.Len(); i++ {
		v := v.Elem().Index(i)

		if v.Kind() == reflect.Ptr && v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = reflect.Indirect(v)
		e := st.Index(i)
		ct.ConvertVal(v, e.Addr())
	}

	return nil
}

func guessMember(obj interface{}) reflect.Type {
	v := reflect.ValueOf(obj)
	t := reflect.Indirect(v).Type()

	if t.Kind() != reflect.Slice {
		panic("obj must a slice Type")
	}

	return t.Elem()
}

func loadOptions(val interface{}, opts []option.OptionFunc) (opt *option.Option) {
	if opt = option.FindDefaultOption(val); opt == nil {
		opt = new(option.Option)
	}

	for _, op := range opts {
		op(opt)
	}

	return opt
}

func loadTypeOptions(typ reflect.Type, opts []option.OptionFunc) (opt *option.Option) {
	if opt = option.FindDefaultTypeOption(typ); opt == nil {
		opt = new(option.Option)
	}

	for _, op := range opts {
		op(opt)
	}

	return opt
}
