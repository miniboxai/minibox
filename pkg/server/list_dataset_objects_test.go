package server

import (
	"context"
	"testing"

	"minibox.ai/minibox/pkg/api/v1/types"
	"minibox.ai/minibox/pkg/utils"
)

func TestListDatasetObjects(t *testing.T) {
	prepareTestDatabase()

	gdb := opentestdb()
	defer gdb.Close()
	client := startRpcServer(gdb)

	ctx := context.Background()

	var id = types.ObjectID{Value: &types.ObjectID_Name{Name: "hysios"}}

	objects, err := client.ListDatasetObjects(ctx,
		&types.ListObjectsRequest{
			Id: &id,
		})
	if err != nil {
		t.Fatalf("List Objects error: %s", err)
	}

	t.Logf("objects: %s", utils.Prettify(objects))
}
