package docker

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	backend, err := NewBackend()
	if err != nil {
		t.Fatalf("creating Docker backend error %s", err)
	}

	reader, err := backend.client.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		t.Fatalf("ImagePull error %s", err)

	}
	io.Copy(os.Stderr, reader)

	resp, err := backend.client.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
		Tty:   true,
	}, nil, nil, "")
	if err != nil {
		t.Fatalf("ContainerCreate error %s", err)
	}

	if _, err := backend.client.ContainerWait(ctx, resp.ID); err != nil {
		t.Fatalf("ContainerWait error %s", err)
	}

	out, err := backend.client.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		t.Fatalf("ContainerLogs error %s", err)
	}

	io.Copy(os.Stderr, out)
}
