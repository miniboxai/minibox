package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Compute struct {
	ID     string
	client *client.Client
	ctx    context.Context
}

func (c *Compute) Wait() (status int64, err error) {
	status, err = c.client.ContainerWait(c.ctx, c.ID)
	return
}

func (c *Compute) Output() (io.ReadCloser, error) {
	return c.client.ContainerLogs(c.ctx, c.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true})

}
