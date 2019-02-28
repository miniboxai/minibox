package server

import (
	"context"
	"testing"

	"minibox.ai/pkg/api/v1/types"
	"minibox.ai/pkg/utils"
)

func TestDeleteProject(t *testing.T) {
	prepareTestDatabase()

	gdb := opentestdb()
	defer gdb.Close()
	client := startRpcServer(gdb)

	ctx := context.Background()

	prj, err := client.DeleteProject(ctx, &types.DeleteProjectRequest{
		ProjectId: "2",
	})
	if err != nil {
		t.Fatalf("Delete Project: %s", err)
	}

	t.Logf("project: %s\n", utils.Prettify(prj))
}

func TestDeleteForceProject(t *testing.T) {
	prepareTestDatabase()

	gdb := opentestdb()
	defer gdb.Close()
	client := startRpcServer(gdb)

	ctx := context.Background()

	prj, err := client.DeleteProject(ctx, &types.DeleteProjectRequest{
		ProjectId: "2",
		Force:     true,
	})
	if err != nil {
		t.Fatalf("Delete Project: %s", err)
	}

	t.Logf("project: %s\n", utils.Prettify(prj))
}
