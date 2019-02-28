package server

import (
	"context"
	"testing"

	"minibox.ai/minibox/pkg/api/v1/types"
)

func TestCreateProject(t *testing.T) {
	prepareTestDatabase()

	gdb := opentestdb()
	defer gdb.Close()
	client := startRpcServer(gdb)

	ctx := context.Background()

	prj, err := client.CreateProject(ctx, &types.CreateProjectRequest{
		Namespace: "testns",
		Name:      "mnist",
	})
	if err != nil {
		t.Fatalf("Create Project: %s", err)
	}
	t.Logf("project: %#v", prj)
	if prj.Name != "mnist" {
		t.Fatal("Create Project failed")
	}
}
