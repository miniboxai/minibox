package server

import (
	"context"

	"minibox.ai/minibox/pkg/api/v1/types"
)

func (s *Server) DeleteProject(ctx context.Context, in *types.DeleteProjectRequest) (*types.Project, error) {
	// curUsr, ok := s.currentUser(ctx)
	// if !ok {
	// 	return nil, ErrInvalidUser
	// }

	var prj = types.Project{
		ID: in.ProjectId,
	}

	if err := s.storage.Delete(&prj); err != nil {
		return nil, err
	}
	return &prj, nil
}
