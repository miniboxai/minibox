package server

import (
	"context"

	"minibox.ai/minibox/pkg/api/v1/types"
	"minibox.ai/minibox/pkg/server/internal/option"
)

func (s *Server) PureDataset(ctx context.Context, in *types.RestoreDatasetRequest) (*types.Dataset, error) {
	var (
		curUsr *types.User
		ok     bool
	)

	// 当前用户, 同时加载组织
	curUsr, ok = s.currentUserWith(ctx, LoadOrgs())
	if !ok {
		return nil, ErrInvalidUser
	}

	dataset := types.Dataset{
		ID: in.DatasetId,
	}

	if err := s.storage.Load(&dataset, option.Undeleted()); err != nil {
		return nil, err
	}

	if dataset.Namespace == curUsr.Namespace {
		goto Exec
	} else if nss := s.orgsNamespaces(curUsr.Organizations); !s.inNamespaces(nss, dataset.Namespace) {
		return nil, ErrInvalidNamespace(dataset.Namespace)
	}

Exec:
	// update dataset
	if err := s.storage.Delete(&dataset, option.Force()); err != nil {
		return nil, err
	}

	return &dataset, nil
}
