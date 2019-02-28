package server

import (
	"bytes"
	"context"
	"io"
	"testing"

	"minibox.ai/pkg/api/v1/types"
)

func TestGetDatasetObject(t *testing.T) {
	prepareTestDatabase()

	gdb := opentestdb()
	defer gdb.Close()
	client := startRpcServer(gdb)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers

	var (
		val = types.ObjectID_Name{"testns"}
		key = types.ObjectID{
			Value: &val,
		}
		blob = "10.txt"
		in   = types.GetObjectRequest{DatasetId: &key, Blob: blob}
	)
	stream, err := client.GetDatasetObject(ctx, &in)
	if err != nil {
		t.Fatalf("GetDatasetObject error %s", err)
	}

	r := NewStreamReader(stream)
	var buf bytes.Buffer
	io.Copy(&buf, r)
	if buf.String() != "hello world" {
		t.Fatalf("contents is error")
	}
	t.Log(buf.String())
}
