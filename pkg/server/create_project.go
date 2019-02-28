package server

import (
	"context"

	"minibox.ai/pkg/api/v1/types"
)

func (s *Server) CreateProject(ctx context.Context, in *types.CreateProjectRequest) (*types.Project, error) {
	curUsr, ok := s.currentUser(ctx)
	if !ok {
		return nil, ErrInvalidUser
	}

	var prj = types.Project{
		Namespace: in.Namespace,
		Name:      in.Name,
		Author:    curUsr,
	}

	if in.Namespace == curUsr.Namespace {
		if err := s.storage.Store(&prj); err != nil {
			return nil, err
		}
		return &prj, nil
	}

	// orgs, err := s.UserOrganizations(ctx, types.UserOrganizationRequest{UserId: curUsr.ID})

	// for _, org := range orgs {
	// 	if org.Namespace == curUsr.Namespace {
	// 		// TODO: check this permission
	// 		//
	// 	}
	// }

	return nil, ErrInvalidNamespace(in.Namespace)
}
