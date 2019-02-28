package server

import (
	"context"
	"testing"

	"minibox.ai/pkg/api/v1/types"
	"minibox.ai/pkg/utils"
)

func TestDeleteDataset(t *testing.T) {
	prepareTestDatabase()

	gdb := opentestdb()
	defer gdb.Close()
	client := startRpcServer(gdb)

	ctx := context.Background()

	ds, err := client.CreateDataset(ctx, &types.CreateDatasetRequest{
		Namespace: "test_org",
		Name:      "mnist",
	})
	if err != nil {
		t.Fatalf("Create Dataset: %s", err)
	}

	dataset, err := client.DeleteDataset(ctx, &types.DeleteDatasetRequest{
		DatasetId: ds.ID,
	})

	if err != nil {
		t.Fatalf("Delete Dataset: %s", err)
	}

	t.Logf("Dataset: %s\n", utils.Prettify(dataset))
}
