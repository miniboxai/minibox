package option

import (
	"reflect"

	"github.com/jinzhu/gorm"
)

type DB = gorm.DB

type FilterHandle struct {
	Handles []FilterFunc
	// Chain   ChainFunc
}

// type ChainFunc func(query interface{}, args ...interface{}) *DB
type FilterFunc func(*DB) *DB

func NewFilters() *FilterHandle {
	return &FilterHandle{}
}

type Map map[string]interface{}

// NewFilters().Map(Map{"user_id": 1})
func (flt *FilterHandle) Map(m Map) *FilterHandle {

	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Where(m)
	})
	return flt
}

func (flt *FilterHandle) OrMap(m Map) *FilterHandle {

	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Or(m)
	})
	return flt
}

func (flt *FilterHandle) NotMap(m Map) *FilterHandle {

	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Not(m)
	})
	return flt
}

// NewFilters().Field("user_id", 1)
func (flt *FilterHandle) Field(name string, val interface{}) *FilterHandle {
	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Where(name+" = ?", val)
	})
	return flt
}

func (flt *FilterHandle) OrField(name string, val interface{}) *FilterHandle {
	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Or(name+" = ?", val)
	})
	return flt
}

func (flt *FilterHandle) NotField(name string, val interface{}) *FilterHandle {
	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Not(name+" = ?", val)
	})
	return flt
}

// NewFilters().FieldOp("user_id", ">", 3)
func (flt *FilterHandle) FieldOp(name string, op string, val interface{}) *FilterHandle {
	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Where(name+" "+op+" ?", val)
	})
	return flt
}

func (flt *FilterHandle) OrFieldOp(name string, op string, val interface{}) *FilterHandle {
	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Or(name+" "+op+" ?", val)
	})
	return flt
}

func (flt *FilterHandle) NotFieldOp(name string, op string, val interface{}) *FilterHandle {
	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Not(name+" "+op+" ?", val)
	})
	return flt
}

// NewFilters().Contains("user_id", []int{1,2,3})
func (flt *FilterHandle) Contains(name string, vals interface{}) *FilterHandle {
	v := reflect.ValueOf(vals)
	if reflect.Indirect(v).Kind() != reflect.Slice {
		panic("val must Slice Type")
	}
	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Where(name+" in (?)", vals)
	})
	return flt
}

func (flt *FilterHandle) OrContains(name string, vals interface{}) *FilterHandle {
	v := reflect.ValueOf(vals)
	if reflect.Indirect(v).Kind() != reflect.Slice {
		panic("val must Slice Type")
	}
	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Or(name+" in (?)", vals)
	})
	return flt
}

func (flt *FilterHandle) NotContains(name string, vals interface{}) *FilterHandle {
	v := reflect.ValueOf(vals)
	if reflect.Indirect(v).Kind() != reflect.Slice {
		panic("val must Slice Type")
	}
	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Not(name+" in (?)", vals)
	})
	return flt
}

// NewFilters().Like("namespace", "%hysios")
func (flt *FilterHandle) Like(name string, query string) *FilterHandle {
	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Where(name+" like ?", query)
	})
	return flt
}

func (flt *FilterHandle) OrLike(name string, query string) *FilterHandle {
	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Or(name+" like ?", query)
	})
	return flt
}

func (flt *FilterHandle) NotLike(name string, query string) *FilterHandle {
	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Not(name+" like ?", query)
	})
	return flt
}

// NewFilters().Where("updated_at > ?", time.Now())
func (flt *FilterHandle) Where(cond string, vals ...interface{}) *FilterHandle {
	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Where(cond, vals...)
	})
	return flt
}

// NewFilters().Where("updated_at > ?", time.Now())
func (flt *FilterHandle) OrWhere(cond string, vals ...interface{}) *FilterHandle {
	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Or(cond, vals...)
	})
	return flt
}

func (flt *FilterHandle) NotWhere(cond string, vals ...interface{}) *FilterHandle {
	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Not(cond, vals...)
	})
	return flt
}

// NewFilters().Struct(&User{ID: 1})
func (flt *FilterHandle) Struct(val interface{}) *FilterHandle {
	v := reflect.ValueOf(val)
	if reflect.Indirect(v).Kind() != reflect.Struct {
		panic("val must Struct Type")
	}

	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Where(val)
	})
	return flt
}

func (flt *FilterHandle) OrStruct(val interface{}) *FilterHandle {
	v := reflect.ValueOf(val)
	if reflect.Indirect(v).Kind() != reflect.Struct {
		panic("val must Struct Type")
	}

	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Or(val)
	})
	return flt
}

func (flt *FilterHandle) NotStruct(val interface{}) *FilterHandle {
	v := reflect.ValueOf(val)
	if reflect.Indirect(v).Kind() != reflect.Struct {
		panic("val must Struct Type")
	}

	flt.Handles = append(flt.Handles, func(db *DB) *DB {
		return db.Not(val)
	})
	return flt
}

func (flt *FilterHandle) Apply(db *DB) *DB {
	for _, hand := range flt.Handles {
		db = hand(db)
	}

	return db
}
