package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gogo/gateway"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	pb "minibox.ai/pkg/api/v1/proto"
)

func TestHttpGateway(t *testing.T) {
	prepareTestDatabase()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	m := new(gateway.JSONPb)
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, m))

	opts := []grpc.DialOption{grpc.WithInsecure()}

	gdb := opentestdb()
	defer gdb.Close()

	s, lis := startServer(gdb)
	pb.RegisterProjectServiceServer(s, NewServer(gdb, nil))

	err := pb.RegisterProjectServiceHandlerFromEndpoint(ctx, mux, lis.Addr().String(), opts)
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(mux)
	defer ts.Close()
	// time.Sleep(1 * time.Second)
	res, err := http.Get(ts.URL + "/api/v1/projects")
	if err != nil {
		log.Fatal(err)
	}
	projects, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", projects)

	// http.ListenAndServe(":11334", mux)

}
