package apiserver

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"

	osin "github.com/RangelReale/osin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"

	"minibox.ai/pkg/server"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

type GRPCServer struct {
	apiServer *ApiServer
}

func ensureAuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")

	if err != nil {
		return nil, err
	}

	tokenInfo, err := parseToken(token)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}
	// grpc_ctxtags.Extract(ctx).Set("auth.sub", userClaimFromToken(tokenInfo))
	newCtx := context.WithValue(ctx, "tokenInfo", tokenInfo)
	if usr, err := getUserByAccessToken(tokenInfo); err != nil {
		if usr, err = getUserFromModel(tokenInfo); err != nil {
			return nil, fmt.Errorf("can't get user by token: %s ", err)
		}
		newCtx = context.WithValue(newCtx, "user", usr)
	} else {
		newCtx = context.WithValue(newCtx, "user", usr)
	}

	return newCtx, nil
}

func parseToken(token string) (*osin.AccessData, error) {
	if len(token) == 0 {
		return nil, errors.New("token string is empty")
	}

	if ad, err := storage.LoadAccess(token); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	} else {
		return ad, nil
	}
}

func userClaimFromToken(ad *osin.AccessData) string {
	if usr, err := getUserByAccessToken(ad); err != nil {
		return "(unknown)"
	} else {
		return usr.Namespace
	}
}

func NewGRPCServer(svr *ApiServer) *GRPCServer {
	return &GRPCServer{
		apiServer: svr,
	}
}

func (gsvr *GRPCServer) Listen(addr string) error {
	log.Printf("server starting on port %s...\n", addr)

	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	s := server.NewServer(db.DB, &server.ServerOption{
		AuthFunc:    ensureAuthFunc,
		Credentials: credentials.NewServerTLSFromCert(&cert),
		SetLogger: func(meth server.GRPCMethod) interface{} {
			if meth == server.Stream {
				return grpc_zap.StreamServerInterceptor(logger)
			} else {
				return grpc_zap.UnaryServerInterceptor(logger)
			}
		},
	})
	s.Registered()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return s.GRPCServer().Serve(lis)
}
