package v1

import (
	"context"
	"net/url"
	"reflect"
	"strings"

	"github.com/kr/pretty"

	gerrors "github.com/pkg/errors"
	"minibox.ai/minibox/pkg/core/errors"
	"minibox.ai/minibox/pkg/core/node"
	"minibox.ai/minibox/pkg/core/option"
	u "minibox.ai/minibox/pkg/core/utils"
	"minibox.ai/minibox/pkg/logger"
)

var prioriEvals = []string{"framework", "datadir", "logdir", "workdir"}

type Map = u.Map
type Any = u.Any

func (cfg *Config) Parse(m map[interface{}]interface{}, opt ...option.ParseOpt) (err error) {
	vm := u.VMap(m)

	defer errorHandle(func(_err error) {
		err = gerrors.WithStack(_err)
		logger.S().Errorf("parse error: %#v\n", err)
	})

	cfg.loadOpts(opt...)
	cfg.nodes = make(map[string]node.Noder)

	train := vm.Map("train")
	logger.S().Debugf("train: % #v", pretty.Formatter(train))

	err = train.RangeSKey(func(key string, val Any) error {
		logger.S().Debugf("config %s: %#v", key, val)
		switch key {
		case "framework":
			fram := cfg.parseFramework(val)
			cv := compileStruct(fram)
			cv.BoundTo(&cfg.Framework)
			cfg.addNode("framework", cv)
			cfg.hasFrame = true
		case "cmd":
			switch v := val.(type) {
			case string:
				cv := compileFunc(v, func(ctnt string, ctx context.Context) interface{} {
					return strings.Split(ctnt, " ")
				})
				cv.BoundTo(&cfg.Cmd)
				cfg.addNode("cmd.start", cv)
			case []interface{}:
				cv := compileStringSlice(v)
				cv.BoundTo(&cfg.Cmd)
				cfg.addNode("cmd.start", cv)
			default:
			}
		case "dataset":
			datasets := u.MustStringSlice(val, "dataset")
			cfg.DataSets = make([]url.URL, 0, len(datasets))
			for _, ds := range datasets {
				if u, err := url.ParseRequestURI(ds); err == nil {
					cfg.DataSets = append(cfg.DataSets, *u)
				}
			}
		case "dirs":
			val := u.Any2Map(val)
			if data, ok := val["data"]; ok {
				cv := compileString(data)
				cv.BoundTo(&cfg.Dirs.Data)
				cfg.addNode("datadir", cv)
			}

			if log, ok := val["log"]; ok {
				cv := compileString(log)
				cv.BoundTo(&cfg.Dirs.Logs)
				cfg.addNode("logdir", cv)
			}
		case "workdir":
			workdir := u.MustString(val, "workdir")
			cv := compileString(workdir)
			cv.BoundTo(&cfg.Workdir)
			cfg.addNode("workdir", cv)
		case "env":
			if v, ok := val.([]interface{}); !ok {
				panic(&errors.ErrInvalidItemType{Name: key, Type: reflect.TypeOf([]string{})})
			} else {
				cv := compileEnvs(v)
				cv.BoundTo(&cfg.Env)
				cfg.addNode("env", cv)
			}
		case "dashboard":
		case "notebook":
		default:
			logger.S().Debugf("cfg.opts %#v", cfg.opts)
			if cfg.opts.Strict {
				return &errors.ErrInvalidConfigItem{Name: key}
			}
		}
		return nil
	})

	if _, ok := cfg.nodes["framework"]; !ok {
		panic(ErrMissingFramework)
	}

	cfg.parseDefaults()
	return err
}

func (cfg *Config) Eval(ctx context.Context) error {
	return Nodes(cfg.nodes).Eval(ctx)
}

func (cfg *Config) addNode(name string, node node.Noder) {
	cfg.nodes[name] = node
}

func (cfg *Config) defaultFrameworkCompile(zero interface{}, handle func(Map, context.Context) interface{}) *node.NodeFunc {
	return compileFunc(zero, func(ctnt string, ctx context.Context) interface{} {
		defaults := cfg.Framework.Defaults()
		return handle(defaults, ctx)
	})
}

func (cfg *Config) parseDefaults() {
	nodes := cfg.nodes

	// 获取默认启动命令
	if _, ok := nodes["cmd.start"]; !ok {
		zero := []string{}
		cv := cfg.defaultFrameworkCompile(zero, func(defaults Map, _ context.Context) interface{} {
			if cmd, ok := defaults["cmd"].([]interface{}); ok {
				return toStringSlice(cmd)
			}
			return zero
		})

		cv.BoundTo(&cfg.Cmd)
		cfg.addNode("cmd.start", cv)
	}

	// 数据目录
	if _, ok := nodes["datadir"]; !ok {
		// zero := []string{}
		cv := cfg.defaultFrameworkCompile("", func(defaults Map, _ context.Context) interface{} {
			return u.VMap(defaults).Map("dirs").String("data")
		})

		cv.BoundTo(&cfg.Dirs.Data)
		cfg.addNode("datadir", cv)
	}

	// 日志目录
	if _, ok := nodes["logdir"]; !ok {
		cv := cfg.defaultFrameworkCompile("", func(defaults Map, _ context.Context) interface{} {
			return u.VMap(defaults).Map("dirs").String("log")
		})

		cv.BoundTo(&cfg.Dirs.Logs)
		cfg.addNode("logdir", cv)
	}

	// 工作目录
	if _, ok := nodes["workdir"]; !ok {
		zero := ""
		cv := cfg.defaultFrameworkCompile(zero, func(defaults Map, ctx context.Context) interface{} {
			if workdir, ok := defaults["workdir"].(string); ok {
				val := compileTmpl(workdir, ctx)
				return val
			}
			return zero
		})
		cv.BoundTo(&cfg.Workdir)
		cfg.addNode("workdir", cv)
	}

	// 环境变量
	if _, ok := nodes["env"]; !ok {
		zero := []Env{}

		cv := cfg.defaultFrameworkCompile(zero, func(defaults Map, ctx context.Context) interface{} {
			ss, ok := defaults["env"].([]interface{})
			if !ok {
				return zero
			}

			ess := toStringSlice(ss)
			var envs = make([]Env, 0, len(ss))
			for _, cenv := range ess {
				val := compileTmpl(cenv, ctx)
				if env, err := ParseEnv(val); err != nil {
					panic(err)
				} else {
					envs = append(envs, *env)
				}
			}
			return envs
		})

		cv.BoundTo(&cfg.Env)
		cfg.addNode("env", cv)
	}
}

func (cfg *Config) parseFramework(val Any) *Framework {
	v, err := ParseFramework(val)
	if err != nil {
		panic(err)
	}
	// cfg.Framework = *v
	return v
}

func minusSlice(origin, toRemove []string) []string {
	for _, rm := range toRemove {
		for i, ori := range origin {
			if ori == rm {
				if i == 0 {
					origin = origin[1:]
				} else {
					origin = append(origin[0:i], origin[i+1:]...)
				}
				break
			}
		}
	}
	return origin
}

func toStringSlice(val []interface{}) []string {
	var ret = make([]string, 0, len(val))
	for _, m := range val {
		if s, ok := m.(string); ok {
			ret = append(ret, s)
		}
	}
	return ret
}

func compileTmpl(tmpl string, ctx context.Context) string {
	var n node.NodeString
	n.Compile(tmpl)
	if v, err := n.Eval(ctx); err != nil {
		panic(err)
	} else {
		return v.String()
	}
}
