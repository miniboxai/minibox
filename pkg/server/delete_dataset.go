package server

import (
	"context"

	"minibox.ai/pkg/api/v1/types"
)

func (s *Server) DeleteDataset(ctx context.Context, in *types.DeleteDatasetRequest) (*types.Dataset, error) {
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

	// load dataset
	if err := s.storage.Load(&dataset); err != nil {
		return nil, err
	}

	if dataset.Namespace == curUsr.Namespace {
		goto Exec
	} else if nss := s.orgsNamespaces(curUsr.Organizations); !s.inNamespaces(nss, dataset.Namespace) {
		return nil, ErrInvalidNamespace(dataset.Namespace)
	}

Exec:
	if err := s.storage.Delete(&dataset); err != nil {
		return nil, err
	}

	return &dataset, nil
}
