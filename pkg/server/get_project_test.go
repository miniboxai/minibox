package server

import (
	"context"
	"testing"

	"minibox.ai/minibox/pkg/api/v1/types"
	"minibox.ai/minibox/pkg/utils"
)

func TestGetProject(t *testing.T) {
	prepareTestDatabase()

	gdb := opentestdb()
	defer gdb.Close()
	client := startRpcServer(gdb)

	ctx := context.Background()

	prj, err := client.GetProject(ctx, &types.GetProjectRequest{
		ProjectId: "2",
	})
	if err != nil {
		t.Fatalf("Get Project: %s", err)
	}

	t.Logf("project: %s\n", utils.Prettify(prj))

	prj, err = client.GetProject(ctx, &types.GetProjectRequest{
		ProjectId: "3",
	})

	if err == nil {
		t.Fatalf("access private must be failed Get Project: %s", err)
	}
	t.Logf("must be error: %s", err)
}
