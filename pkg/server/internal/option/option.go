package option

import (
	"reflect"
	"sync"
)

type Option struct {
	Limit   int           // 限制返回数量
	Offset  int           // 查询结果偏移, 用于处理翻页
	Order   []string      // 排序方式
	Fields  []string      // 字段列表
	Filters *FilterHandle // 过滤条件

	SubLoads    []Loads // 同步加载模块
	Jointable   *JoinLoads
	Force       bool // 强制执行, 用于去除软删除
	Undeleted   bool // 显示删除对象
	AssignAttrs map[string]interface{}
	Selects     []string
	Rejects     []string

	Updates []UpdateField
}

type Loads struct {
	Name string
	Args []interface{}
}

type JoinLoads struct {
	Table interface{}
	Arg   string
}

type UpdateField struct {
	Name  string
	Value reflect.Value
}

type OptionFunc func(*Option) error

func Limit(l int) OptionFunc {
	return func(opt *Option) error {
		opt.Limit = l
		return nil
	}
}

func Offset(o int) OptionFunc {
	return func(opt *Option) error {
		opt.Offset = o
		return nil
	}
}

func Range(oft, lim int) OptionFunc {
	return func(opt *Option) error {
		opt.Offset = oft
		opt.Limit = lim
		return nil
	}
}

func Order(ord string) OptionFunc {
	return func(opt *Option) error {
		opt.Order = append(opt.Order, ord)
		return nil
	}
}

func Fields(fields []string) OptionFunc {
	return func(opt *Option) error {
		opt.Fields = append(opt.Fields, fields...)
		return nil
	}
}

func WithSub(sub string, args ...interface{}) OptionFunc {
	return func(opt *Option) error {
		var load = Loads{
			Name: sub,
			Args: args,
		}

		opt.SubLoads = append(opt.SubLoads, load)
		return nil
	}
}

func Include(sub string, args ...interface{}) OptionFunc {
	return func(opt *Option) error {
		var load = Loads{
			Name: sub,
			Args: args,
		}

		opt.SubLoads = append(opt.SubLoads, load)
		return nil
	}
}

func UpdateWith(field string, arg interface{}) OptionFunc {
	return func(opt *Option) error {
		var upd = UpdateField{
			Name:  field,
			Value: reflect.ValueOf(arg),
		}

		opt.Updates = append(opt.Updates, upd)
		return nil
	}
}

func Jointable(tab interface{}, arg string) OptionFunc {
	return func(opt *Option) error {
		opt.Jointable = &JoinLoads{
			Table: tab,
			Arg:   arg,
		}
		return nil
	}
}

func Force() OptionFunc {
	return func(opt *Option) error {
		opt.Force = true
		return nil
	}
}

func Undeleted() OptionFunc {
	return func(opt *Option) error {
		opt.Undeleted = true
		return nil
	}
}

func Attrs(attrs map[string]interface{}) OptionFunc {
	return func(opt *Option) error {
		opt.AssignAttrs = attrs
		return nil
	}
}

func Select(sel ...string) OptionFunc {
	return func(opt *Option) error {
		opt.Selects = append(opt.Selects, sel...)
		return nil
	}
}

func Reject(rej ...string) OptionFunc {
	return func(opt *Option) error {
		opt.Rejects = append(opt.Rejects, rej...)
		return nil
	}
}

func Filters(hand *FilterHandle) OptionFunc {
	return func(opt *Option) error {
		opt.Filters = hand
		return nil
	}
}

// func join(*db.Scope) *db.Scope {
// }

// func Filter(args interface{}...}) OptionFunc {
// 	return func(opt *Option) error {
// 		opt.Fields = append(opt.Fields, args...)
// 		return nil
// 	}
// }

// func FilterStr(exp string, args interface{}...) OptionFunc {
// 	return func(opt *Option) error {
// 		opt.Fields = append(opt.Fields, args...)
// 		return nil
// 	}
// }

// func FilterMap(m map[string]interface{}) OptionFunc {
// 	return func(opt *Option) error {
// 		opt.Fields = append(opt.Fields, args...)
// 		return nil
// 	}
// }

type filterMap struct {
}

type filterStr struct {
}

type fitlerStruct struct {
}

var (
	mw             sync.RWMutex
	defaultOptions = make(map[reflect.Type]*Option)
)

func RegisterDefault(val interface{}, defOpt *Option) {
	mw.Lock()
	defer mw.Unlock()
	t := reflect.TypeOf(val)
	defaultOptions[t] = defOpt
}

func FindDefaultOption(val interface{}) *Option {
	mw.RLock()
	defer mw.RUnlock()

	t := reflect.TypeOf(val)
	if opt, ok := defaultOptions[t]; ok {
		return opt
	} else {
		return nil
	}
}

func FindDefaultTypeOption(t reflect.Type) *Option {
	mw.RLock()
	defer mw.RUnlock()

	if opt, ok := defaultOptions[t]; ok {
		return opt
	} else {
		return nil
	}
}
