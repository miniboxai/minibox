package server

import (
	"context"

	"minibox.ai/minibox/pkg/api/v1/types"
	"minibox.ai/minibox/pkg/server/internal/option"
)

// 获取项目的详细信息
//
func (s *Server) GetProject(ctx context.Context, in *types.GetProjectRequest) (*types.Project, error) {
	curUsr, ok := s.currentUser(ctx)
	if !ok {
		return nil, ErrInvalidUser
	}

	var prj = types.Project{
		ID: in.ProjectId,
	}

	if err := s.storage.Load(&prj, option.WithSub("User")); err != nil {
		return nil, err
	}

	if prj.Private && prj.Author.ID != curUsr.ID {
		return nil, ErrDenyPrivateProject
	}

	return &prj, nil
}
