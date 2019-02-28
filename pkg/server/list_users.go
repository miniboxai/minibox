package server

import (
	"context"

	"minibox.ai/minibox/pkg/api/v1/types"
	"minibox.ai/minibox/pkg/server/internal/option"
)

func (s *Server) ListUsers(ctx context.Context, in *types.ListUsersRequest) (out *types.UsersReply, err error) {
	var (
		usrs = make([]*types.User, 0, 30)
	)

	if err := s.storage.List(&usrs); err != nil {
		return nil, err
	}

	out = &types.UsersReply{Users: usrs}

	return
}

func init() {
	option.RegisterDefault(&[]*types.User{}, &option.Option{
		Limit: 30,
	})
}
