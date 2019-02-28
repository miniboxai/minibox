package server

import (
	"context"
	"testing"

	"minibox.ai/pkg/api/v1/types"
	"minibox.ai/pkg/utils"
)

func TestListUsers(t *testing.T) {
	prepareTestDatabase()

	gdb := opentestdb()
	defer gdb.Close()
	client := startRpcServer(gdb)

	ctx := context.Background()

	usrs, err := client.ListUsers(ctx, &types.ListUsersRequest{})
	if err != nil {
		t.Fatalf("List Users error: %s", err)
	}

	t.Logf("List Users: %s\n", utils.Prettify(usrs))
}
