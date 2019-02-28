package server

import (
	"context"
	"testing"

	"minibox.ai/minibox/pkg/api/v1/types"
)

func TestUpdateProject(t *testing.T) {
	prepareTestDatabase()

	gdb := opentestdb()
	defer gdb.Close()
	client := startRpcServer(gdb)

	ctx := context.Background()

	prj, err := client.UpdateProject(ctx, &types.UpdateProjectRequest{
		ProjectId: "1",
		Name:      "new_name",
	})

	if err != nil {
		t.Fatalf("Update Project error: %s", err)
	}
	t.Logf("Project: %#v", prj)

	update := &types.UpdateProjectRequest{
		ProjectId: "1",
		Name:      "new_name",
		Private:   &types.UpdateProjectRequest_Value{false},
	}

	prj, err = client.UpdateProject(ctx, update)

	if err != nil {
		t.Fatalf("Update Project error: %s", err)
	}
	t.Logf("Project: %#v", prj)
}
