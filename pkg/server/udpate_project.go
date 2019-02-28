package server

import (
	"context"

	"minibox.ai/minibox/pkg/api/v1/types"
	"minibox.ai/minibox/pkg/server/internal/option"
)

func (s *Server) UpdateProject(ctx context.Context, in *types.UpdateProjectRequest) (*types.Project, error) {
	// curUsr, ok := s.currentUser(ctx)
	// if !ok {
	// 	return nil, ErrInvalidUser
	// }

	var prj = types.Project{
		ID: in.ProjectId,
	}

	var attrs = make(map[string]interface{})

	if in.GetPrivate() != nil {
		attrs["private"] = in.GetPrivate()
	}

	if len(in.Name) > 0 {
		attrs["name"] = in.Name
	}

	if err := s.storage.Store(&prj, option.Attrs(attrs)); err != nil {
		return nil, err
	}
	return &prj, nil
}
