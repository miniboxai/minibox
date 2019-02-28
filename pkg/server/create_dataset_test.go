package server

import (
	"context"
	"testing"

	"minibox.ai/pkg/api/v1/types"
	"minibox.ai/pkg/utils"
)

func TestCreateDataset(t *testing.T) {
	prepareTestDatabase()

	gdb := opentestdb()
	defer gdb.Close()
	client := startRpcServer(gdb)

	ctx := context.Background()

	dataset, err := client.CreateDataset(ctx, &types.CreateDatasetRequest{
		Name: "mnist",
	})
	if err != nil {
		t.Fatalf("Create Dataset: %s", err)
	}

	t.Logf("Dataset: %s\n", utils.Prettify(dataset))
	if dataset.Name != "mnist" {
		t.Fatal("Create Dataset failed")
	}

	dataset, err = client.CreateDataset(ctx, &types.CreateDatasetRequest{
		Name:    "mnist2",
		Private: true,
	})
	if err != nil {
		t.Fatalf("Create Private Dataset: %s", err)
	}

	t.Logf("Private Dataset: %s\n", utils.Prettify(dataset))
	if !dataset.Private {
		t.Fatal("Create Private Dataset failed")
	}

	dataset, err = client.CreateDataset(ctx, &types.CreateDatasetRequest{
		Name:      "nametest",
		Namespace: "test4",
	})
	if err == nil {
		t.Fatalf("Create non-exsity `test4` Dataset must failed: %s", err)
	}

	dataset, err = client.CreateDataset(ctx, &types.CreateDatasetRequest{
		Name:      "nametest",
		Namespace: "test_org",
	})

	if err != nil {
		t.Fatalf("Create Private Dataset: %s", err)
	}

	t.Logf("Create another namespace Dataset: %s\n", utils.Prettify(dataset))

}
