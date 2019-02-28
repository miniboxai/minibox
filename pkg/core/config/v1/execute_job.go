package v1

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"path/filepath"

	"minibox.ai/minibox/pkg/core/job"
	"minibox.ai/minibox/pkg/core/option"
	"minibox.ai/minibox/pkg/core/volume"
)

func (cfg *Config) Start(ctxt context.Context, opts *option.StartJobOption) (*job.Job, error) {
	var (
		j   job.Job
		opt = &option.StartJobOption{
			Name:     "",
			RootPath: "",
		}
	)

	if opts != nil {
		opt = opts
	}

	j.Name = opt.Name
	j.RootPath = opt.RootPath

	if err := cfg.Eval(ctxt); err != nil {
		return nil, err
	}

	var volumes = []job.VolumeMap{
		job.VolumeMap{
			Src:  &volume.ProjectVolume{j.Name, j.RootPath},
			Dest: cfg.Workdir,
		},
	}

	tWorker := &job.TaskSpec{
		Labels: []job.Label{
			job.Label{"task.kind", "worker"},
			job.Label{"mini.app", "project"},
		},
		Name:       "job-master",
		Image:      cfg.Framework.Image(),
		Volumes:    volumes,
		WorkingDir: cfg.Workdir,
		Cmd:        cfg.Cmd,
		Env:        envToStringSlice(cfg.Env),
	}

	// Name     string
	// Image    string
	// Ports    []PortMap
	// Volumnes []VolumeMap
	// Links    []*TaskSpec
	// Cmd      []string
	// master Pod
	// parameter Pod

	// data volumes Pod
	for _, ds := range cfg.DataSets {
		log.Printf("dataset: %#v", ds)
		source, pat, _ := parseDataset(ds)
		switch source {
		case HubData, UserData:
			task := cfg.buildHubDataTask(&j, pat)
			j.Tasks = append(j.Tasks, task)
		case WideData:
		case LocalData:
		default:
			log.Printf("invalid datset source: %s\n", ds)
		}
	}

	// dashboard
	// notebook

	j.Tasks = append(j.Tasks, tWorker)
	return &j, nil
}

const prepareDataImage = "predata"

func (cfg *Config) buildHubDataTask(j *job.Job, pat string) *job.TaskSpec {
	// var dataDir = cfg.Dirs.Data

	task := &job.TaskSpec{
		Labels: []job.Label{
			job.Label{"task.kind", "prepare-data"},
			job.Label{"mini.app", "project"},
		},
		Name:       fmt.Sprintf("%s-prepare-data", j.Name),
		Image:      prepareDataImage,
		Volumes:    make([]job.VolumeMap, 0),
		WorkingDir: cfg.Workdir,
		Cmd:        []string{pat},
		Env:        envToStringSlice(cfg.Env),
	}
	return task
}

type DataSource int

const (
	HubData DataSource = iota
	UserData
	PrivateData
	WideData
	LocalData
	InvalidData
)

var hosthub = "minibox.ai"

func isLocalPath(pat string) bool {
	if len(pat) == 0 {
		return false
	}

	if []rune(pat)[0] == '.' {
		return true
	} else if len(filepath.VolumeName(pat)) > 0 {
		return true
	}

	return false
}

func parseDataset(u url.URL) (source DataSource, pat, resURI string) {
	dir, file := filepath.Split(u.Path)

	if empty(u.Host) && empty(dir) && len(file) > 0 {
		source = HubData
		pat = file
		resURI = fmt.Sprintf("https://%s/datasets/%s", hosthub, file)
	} else if empty(u.Host) && len(dir) > 0 && len(file) > 0 {
		source = UserData
		pat = u.Path
		resURI = fmt.Sprintf("https://%s/%s/datasets/%s", hosthub, dir, file)
	} else if len(u.Host) > 0 && u.Host != hosthub {
		source = WideData
		pat = u.String()
		resURI = pat
	} else if u.Scheme == "file" || isLocalPath(u.Path) {
		source = LocalData
		pat = u.Path
		resURI = pat
	} else {
		source = InvalidData
	}

	return
}

func envToStringSlice(env []Env) []string {
	ss := make([]string, 0, len(env))
	for _, e := range env {
		ss = append(ss, fmt.Sprintf("%s=%s", e.Name, e.Value))
	}
	return ss
}
