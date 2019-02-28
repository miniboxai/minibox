package server

import (
	"context"
	"testing"

	"minibox.ai/minibox/pkg/api/v1/types"
)

func TestListDatasets(t *testing.T) {
	prepareTestDatabase()

	gdb := opentestdb()
	defer gdb.Close()
	client := startRpcServer(gdb)

	ctx := context.Background()

	datasets, err := client.ListDatasets(ctx,
		&types.ListDatasetsRequest{
			Namespace:   "testns",
			ShowPrivate: true,
		})
	if err != nil {
		t.Fatalf("List Datasets error: %s", err)
	}

	if len(datasets.Datasets) != 2 {
		t.Fatalf("List Datasets error: %s", err)
	}

	datasets, err = client.ListDatasets(ctx,
		&types.ListDatasetsRequest{
			Namespace: "testns",
		})
	if err != nil {
		t.Fatalf("List Datasets error: %s", err)
	}

	if len(datasets.Datasets) != 1 {
		t.Fatalf("List Datasets error: %s", err)
	}
}
