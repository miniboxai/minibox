package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"minibox.ai/pkg/backend/option"
	"minibox.ai/pkg/core/job"

	"minibox.ai/pkg/logger"

	"minibox.ai/pkg/utils"
	"minibox.ai/pkg/utils/pretty"
)

type Backend struct {
	client   *client.Client
	networks []types.NetworkResource
}

type Container struct {
	ID string
}

func NewBackend() (*Backend, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return &Backend{
		client: cli,
	}, nil
}

func (be *Backend) Execute(j *job.Job, ctx context.Context, opt *option.ExecuteOption) (err error) {
	var handle = func(h func(error) error) {
		r := recover()
		if rerr, ok := r.(error); ok {
			err = h(rerr)
		}
	}

	defer handle(func(err error) error {
		logger.S().Errorf("Execute error: %s", err)
		return err
	})
	// var bjob Job
	j.ID = utils.GenerateRandomID()
	wrk := be.lookupTask("worker", j)
	be.prepareWorker(ctx, wrk)


	if err := be.prepareImage(ctx, j, wrk); err != nil {
		return err
	}

	comput, err := be.createCompute(ctx, j, wrk)
	if err != nil {
		return err
	}
	// resp, err := be.client.ContainerCreate(ctx, &container.Config{
	// 	WorkingDir:   tWorker.WorkingDir,
	// 	Image:        tWorker.Image,
	// 	Cmd:          tWorker.Cmd,
	// 	Env:          tWorker.Env,
	// 	Tty:          true,
	// 	AttachStdin:  true,
	// 	AttachStdout: true, // Attach the standard output
	// 	AttachStderr: true, // Attach the standard error
	// }, &container.HostConfig{
	// 	Mounts: []mount.Mount{
	// 		{
	// 			Type:   mount.TypeBind,
	// 			Source: j.RootPath,
	// 			Target: tWorker.WorkingDir,
	// 		},
	// 		{
	// 			Type:   mount.TypeBind,
	// 			Source: path.Join(j.RootPath, "/notebooks"),
	// 			Target: "/notebooks",
	// 		},
	// 	},
	// }, nil, "")

	// if err != nil {
	// 	panic(fmt.Errorf("ContainerCreate error %s", err))
	// }

	// logger.S().Infof("container: %# v", resp)

	// // status, err := be.client.ContainerWait(ctx, resp.ID)
	// // if err != nil {
	// // 	panic(fmt.Errorf("ContainerWait error %s", err))
	// // }

	// fmt.Printf("Create Job: %s\n", j.ID)

	// if err := be.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Starting Job: %s\n", j.ID)

	if opt.WaitStatus {

		go func() {
			var ch = make(chan int)
			j.SetState(job.StateExit, ch)

			status, err := comput.Wait()
			if err != nil {
				panic(fmt.Errorf("ContainerWait error %s", err))
			}

			logger.S().Infof("Container State: %d", status)
			ch <- int(status)
		}()
	}

	if opt.Log {
		out, err := comput.Output()
		if err != nil {
			logger.S().Errorf("ContainerLogs error %s", err)
		}

		j.SetOutput(out)
	}
	// io.Copy(os.Stderr, out)

	return nil
}

func (be *Backend) lookupTaskByLabels(name, value string, j *job.Job) []*job.TaskSpec {
	var tasks []*job.TaskSpec
	for _, task := range j.Tasks {
		for _, label := range task.Labels {
			if label.Name == name && label.Value == value {
				tasks = append(tasks, task)
			}
		}
	}
	return tasks
}

func (be *Backend) lookupTasks(name string, j *job.Job) []*job.TaskSpec {
	return be.lookupTaskByLabels("task.kind", name, j)
}

func (be *Backend) lookupTask(name string, j *job.Job) *job.TaskSpec {
	ts := be.lookupTasks(name, j)
	if len(ts) > 0 {
		return ts[0]
	}

	return nil
}

func (be *Backend) prepareImage(ctx context.Context, j *job.Job, task *job.TaskSpec) error {
	opt := types.ImagePullOptions{}
	rd, err := be.client.ImagePull(ctx, task.Image, opt)

	io.Copy(os.Stderr, rd)
	if err != nil {
		return err
	}
	return nil
}

func (be *Backend) createCompute(ctx context.Context, j *job.Job, task *job.TaskSpec) (*Compute, error) {
	var (
		compute Compute
		mounts  []mount.Mount
	)

	for _, vol := range task.Volumes {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: vol.Src.Source(),
			Target: vol.Dest,
		})
	}

	// []mount.Mount{
	// 	{
	// 		Type:   mount.TypeBind,
	// 		Source: j.RootPath,
	// 		Target: tWorker.WorkingDir,
	// 	},
	// 	{
	// 		Type:   mount.TypeBind,
	// 		Source: path.Join(j.RootPath, "/notebooks"),
	// 		Target: "/notebooks",
	// 	},

	resp, err := be.client.ContainerCreate(ctx, &container.Config{
		WorkingDir:   task.WorkingDir,
		Image:        task.Image,
		Cmd:          task.Cmd,
		Env:          task.Env,
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true, // Attach the standard output
		AttachStderr: true, // Attach the standard error
	}, &container.HostConfig{
		Mounts: mounts,
	}, nil, "")

	if err != nil {
		return nil, fmt.Errorf("ContainerCreate error %s", err)
	}

	compute.ID = resp.ID
	compute.client = be.client
	compute.ctx = ctx
	logger.S().Infof("container: %# v", resp)

	// status, err := be.client.ContainerWait(ctx, resp.ID)
	// if err != nil {
	// 	panic(fmt.Errorf("ContainerWait error %s", err))
	// }

	fmt.Printf("Create Job: %s\n", shortName(j.ID))

	if err := be.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	fmt.Printf("Starting Job: %s\n", shortName(j.ID))

	return &compute, nil
}

func (be *Backend) prepareWorker(ctx context.Context, task *job.TaskSpec) {
	// netres, err := be.Networks(ctx)
	args := filters.NewArgs()
	args.Add("name", "bridge")
	args.Add("driver", "bridge")
	args.Add("scope", "local")

	netres, err := be.FilterNetworks(ctx, types.NetworkListOptions{args})
	if err != nil {
		panic(err)
	}
	logger.S().Infof("networks: % #v", pretty.Formatter(netres))
	be.networks = netres
}

func (be *Backend) Networks(ctx context.Context) ([]types.NetworkResource, error) {
	if netres, err := be.client.NetworkList(ctx, types.NetworkListOptions{}); err != nil {
		return nil, err
	} else {
		return netres, nil
	}
	// func (cli *Client) NetworkList(ctx context.Context, options types.NetworkListOptions) ([]types.NetworkResource, error) {
}

func (be *Backend) FilterNetworks(ctx context.Context, filters types.NetworkListOptions) ([]types.NetworkResource, error) {
	if netres, err := be.client.NetworkList(ctx, filters); err != nil {
		return nil, err
	} else {
		return netres, nil
	}
}
