package server

import (
	"context"

	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	"minibox.ai/minibox/pkg/api/v1/types"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"

	pb "minibox.ai/minibox/pkg/api/v1/proto"

	"minibox.ai/minibox/pkg/models"
	"minibox.ai/minibox/pkg/server/internal/object_store"
	"minibox.ai/minibox/pkg/server/internal/storage"
)

type Server struct {
	grpcserv *grpc.Server
	storage  *storage.Storage
	objects  *object_store.ObjectStore
}

// cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
// if err != nil {
// 	log.Fatalf("failed to load key pair: %s", err)
// }
// opts := []grpc.ServerOption{
// 	grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
// 		grpc_auth.StreamServerInterceptor(ensureAuthFunc),
// 		grpc_zap.StreamServerInterceptor(logger),
// 		grpc_recovery.StreamServerInterceptor(),
// 	)),
// 	// The following grpc.ServerOption adds an interceptor for all unary
// 	// RPCs. To configure an interceptor for streaming RPCs, see:
// 	// https://godoc.org/google.golang.org/grpc#StreamInterceptor
// 	grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
// 		grpc_auth.UnaryServerInterceptor(ensureAuthFunc),
// 		grpc_zap.UnaryServerInterceptor(logger),
// 		grpc_recovery.UnaryServerInterceptor(),
// 	)),
// 	// Enable TLS for all incoming connections.
// 	grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
// }

func NewServer(db *gorm.DB, opt *ServerOption) *Server {
	var (
		server             Server
		streamInterceptors []grpc.StreamServerInterceptor
		unaryInterceptors  []grpc.UnaryServerInterceptor
	)

	if opt.AuthFunc != nil {
		streamInterceptors = append(streamInterceptors,
			grpc_auth.StreamServerInterceptor(opt.AuthFunc))
		unaryInterceptors = append(unaryInterceptors,
			grpc_auth.UnaryServerInterceptor(opt.AuthFunc))
	}

	if opt.Validate {
		streamInterceptors = append(streamInterceptors,
			grpc_validator.StreamServerInterceptor())
		unaryInterceptors = append(unaryInterceptors,
			grpc_validator.UnaryServerInterceptor())
	}

	if opt.SetLogger != nil {
		if interceptor, ok := opt.SetLogger(0).(grpc.StreamServerInterceptor); ok {
			streamInterceptors = append(streamInterceptors,
				interceptor)
		}

		if interceptor, ok := opt.SetLogger(1).(grpc.UnaryServerInterceptor); ok {
			unaryInterceptors = append(unaryInterceptors,
				interceptor)
		}
	}

	if opt.Recovery {
		streamInterceptors = append(streamInterceptors,
			grpc_recovery.StreamServerInterceptor())
		unaryInterceptors = append(unaryInterceptors,
			grpc_recovery.UnaryServerInterceptor())
	}

	grpcOpts := []grpc.ServerOption{
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			streamInterceptors...,
		)),
		// The following grpc.ServerOption adds an interceptor for all unary
		// RPCs. To configure an interceptor for streaming RPCs, see:
		// https://godoc.org/google.golang.org/grpc#StreamInterceptor
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			unaryInterceptors...,
		)),
	}

	if opt.Credentials != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(opt.Credentials))
	}

	server.initStorage(db)
	server.initGrpcserv(grpcOpts)
	server.initObjectStore(&opt.ObjectStore)

	return &server
}

func (s *Server) initStorage(db *gorm.DB) {
	s.storage = storage.NewStorage(db)
}

func (s *Server) initGrpcserv(opts []grpc.ServerOption) {
	s.grpcserv = grpc.NewServer(opts...)
}

func (s *Server) initObjectStore(opt *ObjectStoreOption) {
	var (
		opts []object_store.OptionFunc
	)

	if len(opt.Region) > 0 {
		opts = append(opts, object_store.Region(opt.Region))
	}
	if len(opt.Bucket) > 0 {
		opts = append(opts, object_store.StoreBucket(opt.Bucket))
	}
	if len(opt.Root) > 0 {
		opts = append(opts, object_store.Region(opt.Root))
	}

	s.objects = object_store.NewObjectStore(opt.Key, opt.Secret, opts...)
}

func (s *Server) Registered() {
	pb.RegisterProjectServiceServer(s.grpcserv, s)
	pb.RegisterUserServiceServer(s.grpcserv, s)
	pb.RegisterDatasetServiceServer(s.grpcserv, s)
}

func (s *Server) GRPCServer() *grpc.Server {
	return s.grpcserv
}

func (s *Server) currentUser(ctx context.Context) (*types.User, bool) {
	if curUsr, ok := ctx.Value("user").(*models.User); ok {
		return usrModel2Type(curUsr), true
	} else {
		return nil, false
	}
}

func (s *Server) currentUserWith(ctx context.Context, opts ...UserOption) (*types.User, bool) {
	if curUsr, ok := ctx.Value("user").(*models.User); ok {
		for _, op := range opts {
			op(s, usrModel2Type(curUsr))
		}
		return usrModel2Type(curUsr), true
	} else {
		return nil, false
	}
}

func (s *Server) ListUsersByRole(context.Context, *types.UserRole) (*types.UsersReply, error) {
	panic("not implemented")
}

func (s *Server) UpdateUser(context.Context, *types.UpdateUserRequest) (*types.User, error) {
	panic("not implemented")
}

func (s *Server) DeleteDatasetObject(context.Context, *types.DeleteObjectRequest) (*types.ObjectReply, error) {
	panic("not implemented")
}
