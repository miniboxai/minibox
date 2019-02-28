package server

import (
	"context"
	"testing"

	"minibox.ai/minibox/pkg/api/v1/types"
	"minibox.ai/minibox/pkg/utils"
)

func TestListProject(t *testing.T) {
	prepareTestDatabase()

	gdb := opentestdb()
	defer gdb.Close()
	client := startRpcServer(gdb)

	ctx := context.Background()

	prjs, err := client.ListProject(ctx, &types.User{ID: 2})
	if err != nil {
		t.Fatalf("List Project error: %s", err)
	}

	t.Logf("List project: %s\n", utils.Prettify(prjs))
}

func BenchmarkListProject(b *testing.B) {
	prepareTestDatabase()

	gdb := opentestdb()
	defer gdb.Close()
	client := startRpcServer(gdb)
	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := client.ListProject(ctx, &types.User{ID: 2})
		if err != nil {
			b.Fatalf("List Project error: %s", err)
		}
	}
}

func BenchmarkListProjectParallel(b *testing.B) {
	prepareTestDatabase()

	gdb := opentestdb()
	defer gdb.Close()
	client := startRpcServer(gdb)
	ctx := context.Background()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.ListProject(ctx, &types.User{ID: 2})
			if err != nil {
				b.Fatalf("List Project error: %s", err)
			}
		}
	})
}
