package core

import (
	"fmt"
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
	"minibox.ai/pkg/core/config/v1"
	"minibox.ai/pkg/core/errors"
	"minibox.ai/pkg/core/option"
	"minibox.ai/pkg/core/project"
)

type Config struct {
	Version string
	Train   interface{}
}

func FromConfig(r io.Reader, opts *option.ConfigOpt) (*project.Project, error) {
	var prj project.Project

	if opts == nil {
		return nil, errors.ErrOptionCantNull
	}

	prj = project.Project{
		Name:     opts.ProjectName,
		RootPath: opts.RootPath,
	}

	m := make(map[interface{}]interface{})
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// defaultProject(&prj)

	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		return nil, err
	}
	cfgStruct, err := parseVersion(m)
	if err != nil {
		return nil, err
	}
	cfg, ok := cfgStruct.(project.Parser)
	if !ok {
		return nil, fmt.Errorf("cfgStruct must have Configer interface{}")
	}

	prj.SetConfig(cfg)
	// var initialParams = map[string]interface{}{
	// 	"projectName": prj.Name,
	// }

	if err = cfg.Parse(m); err != nil {
		return nil, err
	}

	// ctx := context.WithValue(
	// 	context.Background(),
	// 	pcontext.ContextKey,
	// 	pcontext.NewContext(initialParams))

	// cfg.Eval(ctx)

	if trainer, ok := cfgStruct.(project.Trainer); !ok {
		return nil, errors.ErrTrainerInterface
	} else {
		prj.SetTrainer(trainer)
		// prj.trainer = trainer
	}
	return &prj, nil
}

func init() {
	// Registry(defaultVersion, &ConfigV1{})
	Registry(defaultVersion, &v1.Config{})
}
