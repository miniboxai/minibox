package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"testing"

	"github.com/go-testfixtures/testfixtures"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"google.golang.org/grpc"
	v1 "minibox.ai/minibox/pkg/api/v1"
	pb "minibox.ai/minibox/pkg/api/v1/proto"
	"minibox.ai/minibox/pkg/api/v1/types"
	"minibox.ai/minibox/pkg/models"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
)

var (
	db       *sql.DB
	fixtures *testfixtures.Context
)

func TestMain(m *testing.M) {
	var err error
	gdb, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("Open sqlite3 failed: %v", err)
	}
	gdb.AutoMigrate(&models.Project{})
	gdb.AutoMigrate(&models.User{})
	gdb.AutoMigrate(&models.Dataset{})
	gdb.AutoMigrate(&models.Organization{})
	gdb.AutoMigrate(&models.UserOrganization{})

	gdb.Close()
	// Open connection with the test database.
	// Do NOT import fixtures in a production database!
	// Existing data would be deleted
	db, err = sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}

	root := os.Getenv("ROOT")
	dir := path.Join(root, "db/fixtures/")
	// creating the context that hold the fixtures
	// see about all compatible databases in this page below
	fixtures, err = testfixtures.NewFolder(db, &testfixtures.SQLite{}, dir)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

func dialGrpc(lis net.Listener) *grpc.ClientConn {
	addr := fmt.Sprintf("%s", lis.Addr())
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}

func newClient(lis net.Listener) *v1.Clients {
	return v1.NewClients(dialGrpc(lis))
}

func opentestdb() *gorm.DB {
	gdb, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("Open sqlite3 failed: %v", err)
	}

	if len(os.Getenv("DEBUG")) > 0 {
		gdb.LogMode(true)
	}

	return gdb
}

func authFunc(ctx context.Context) (context.Context, error) {
	var usr = types.User{
		ID:        2,
		Name:      "test",
		Namespace: "testns",
	}
	ctx = context.WithValue(ctx, "user", &usr)
	return ctx, nil
}

func startServer(db *gorm.DB) (*grpc.Server, net.Listener) {
	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_auth.StreamServerInterceptor(authFunc),
			grpc_validator.StreamServerInterceptor(),
		)),
		// The following grpc.ServerOption adds an interceptor for all unary
		// RPCs. To configure an interceptor for streaming RPCs, see:
		// https://godoc.org/google.golang.org/grpc#StreamInterceptor
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_auth.UnaryServerInterceptor(authFunc),
			grpc_validator.UnaryServerInterceptor(),
		)),
	}

	s := grpc.NewServer(opts...)
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go s.Serve(lis)
	return s, lis
}

func startRpcServer(db *gorm.DB) *v1.Clients {

	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_auth.StreamServerInterceptor(authFunc),
			grpc_validator.StreamServerInterceptor(),
		)),
		// The following grpc.ServerOption adds an interceptor for all unary
		// RPCs. To configure an interceptor for streaming RPCs, see:
		// https://godoc.org/google.golang.org/grpc#StreamInterceptor
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_auth.UnaryServerInterceptor(authFunc),
			grpc_validator.UnaryServerInterceptor(),
		)),
	}

	var (
		appid  = os.Getenv("AWSAPPID")
		secret = os.Getenv("AWSSECRET")
		region = os.Getenv("AWSREGION")
		bucket = os.Getenv("AWSBUCKET")
		root   = "/test/datasets"
	)

	s := grpc.NewServer(opts...)
	ss := NewServer(db, &ServerOption{
		Key:    appid,
		Secret: secret,
		Region: region,
		Bucket: bucket,
		Root:   root,
	})
	pb.RegisterProjectServiceServer(s, ss)
	pb.RegisterUserServiceServer(s, ss)
	pb.RegisterDatasetServiceServer(s, ss)

	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	client := newClient(lis)
	go s.Serve(lis)
	return client
}
