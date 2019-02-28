package server

import (
	"context"

	"minibox.ai/minibox/pkg/api/v1/types"
	"minibox.ai/minibox/pkg/server/internal/option"
)

func (s *Server) ListProjects(ctx context.Context, in *types.ListProjectsRequest) (*types.ProjectsReply, error) {
	var (
		prjs = make([]*types.Project, 0, 30)
	)

	if err := s.storage.List(&prjs, option.WithSub("User")); err != nil {
		return nil, err
	}

	var projs = types.ProjectsReply{Projects: prjs}

	return &projs, nil
}

func init() {
	option.RegisterDefault(&[]*types.Project{}, &option.Option{
		Limit: 30,
	})
}
