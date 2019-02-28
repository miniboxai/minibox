package project

import (
	"context"
	"errors"

	"minibox.ai/pkg/core/job"
	"minibox.ai/pkg/core/option"
)

type Parser interface {
	Parse(map[interface{}]interface{}, ...option.ParseOpt) error
	Eval(context.Context) error
}

type Trainer interface {
	Start(context.Context, *option.StartJobOption) (*job.Job, error)
}

// type Executor interface {
// 	Execute(*job.Job, context.Context, *option.ExecuteOption) error
// }

type Project struct {
	// config *Config
	Name     string
	RootPath string
	config   interface{}
	trainer  Trainer
	// backend  backend.Executor
	// ImageName string
}

func (prj *Project) Start(ctx context.Context, opt *option.StartProjectOption) (*job.Job, error) {
	if opt == nil {
		return nil, errors.New("opt can't be null")
	}

	j, err := prj.trainer.Start(ctx, &option.StartJobOption{
		Name:     opt.ProjectName,
		RootPath: prj.RootPath,
	})

	// j.Project = prj
	if err != nil {
		return nil, err
	}

	// if prj.backend == nil {
	// 	return nil, errors.New("dont have Executor")
	// }

	// var execOpt backend.ExecuteOption
	// if opt.Wait {
	// 	execOpt.WaitStatus = true
	// 	execOpt.Log = true
	// }

	// if err = prj.backend.Execute(j, ctx, &execOpt); err != nil {
	// 	return nil, fmt.Errorf("execute job error %s", err)
	// }

	return j, nil
}

func (prj *Project) SetConfig(cfg interface{}) {
	prj.config = cfg
}

func (prj *Project) SetTrainer(trainer Trainer) {
	prj.trainer = trainer
}

// func (prj *Project) SetBackend(b backend.Executor) {
// 	prj.backend = b
// }
