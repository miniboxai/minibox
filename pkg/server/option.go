package server

import (
	"context"

	"google.golang.org/grpc/credentials"
	"minibox.ai/minibox/pkg/api/v1/types"
	"minibox.ai/minibox/pkg/server/internal/option"
)

type UserOption func(*Server, *types.User)
type GRPCMethod int

const (
	Stream GRPCMethod = iota
	Unary
)

func LoadGroups() UserOption {
	return func(s *Server, usr *types.User) {

	}
}

func LoadOrgs() UserOption {
	return func(s *Server, usr *types.User) {
		var (
			orgs    = make([]*types.Organization, 0, 10)
			filters = option.
				NewFilters().
				Field("owner_id", usr.ID)
		)

		s.storage.List(&orgs, option.Filters(filters))
		usr.Organizations = orgs
	}
}

type ServerOption struct {
	ObjectStore ObjectStoreOption
	AuthFunc    func(ctx context.Context) (context.Context, error)
	SetLogger   func(meth GRPCMethod) interface{}
	Credentials credentials.TransportCredentials
	Recovery    bool
	Log         bool
	Validate    bool
}

type ObjectStoreOption struct {
	Key    string
	Secret string
	Bucket string
	Region string
	Root   string
}
