package server

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/kr/pretty"
)

func TestPutDatasetObject(t *testing.T) {
	prepareTestDatabase()

	gdb := opentestdb()
	defer gdb.Close()
	client := startRpcServer(gdb)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	rd := strings.NewReader("hello world")

	stream, err := client.PutDatasetObject(ctx)
	if err != nil {
		t.Fatalf("PutDatasetObject %s", err)
	}

	w := NewStreamWriter("10.txt", stream)
	_, err = io.Copy(w, rd)
	if err != nil {
		t.Fatalf("Copy Error %s", err)
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("CloseAndRecv Error %s", err)
	}
	t.Logf("reply: %# v", pretty.Formatter(reply))
	// s.Close()
	t.Log("susccess")
}
