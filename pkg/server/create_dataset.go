package server

import (
	"context"
	"strings"

	"minibox.ai/minibox/pkg/api/v1/types"
	"minibox.ai/minibox/pkg/server/internal/option"
	"minibox.ai/minibox/pkg/utils"
)

func (s *Server) orgsNamespaces(orgs []*types.Organization) []string {
	var nss = make([]string, len(orgs))
	for i, org := range orgs {
		nss[i] = org.Namespace
	}
	return nss
}

func (s *Server) inNamespaces(nss []string, namespace string) bool {
	for _, ns := range nss {
		if strings.ToLower(namespace) == strings.ToLower(ns) {
			return true
		}
	}

	return false
}

func (s *Server) CreateDataset(ctx context.Context, in *types.CreateDatasetRequest) (*types.Dataset, error) {
	var (
		curUsr  *types.User
		ok      bool
		dataset types.Dataset
	)

	curUsr, ok = s.currentUserWith(ctx, LoadOrgs())
	if !ok {
		return nil, ErrInvalidUser
	}

	dataset.Name = in.Name

	nss := s.orgsNamespaces(curUsr.Organizations)
	// namespace checking
	if utils.Empty(in.Namespace) {
		in.Namespace = curUsr.Namespace
	} else if in.Namespace == curUsr.Namespace {
		goto Next
	} else if !s.inNamespaces(nss, in.Namespace) {
		return nil, ErrInvalidNamespace(in.Namespace)
	}

Next:
	dataset.Namespace = in.Namespace

	if in.Private {
		// TODO: can create private dataset
		dataset.Private = in.Private
	}
	dataset.Description = in.Description
	dataset.Preprocessing = in.Preprocessing
	if len(in.Format) > 0 {
		// TODO: how many formats
		dataset.Format = in.Format
	}

	if len(in.Resources) > 0 {
		// TODO: management many Resources
		dataset.Resources = in.Resources
	}

	if len(in.Lisense) > 0 {
		dataset.Lisense = in.Lisense
	}

	if len(in.Reference) > 0 {
		// TODO: management many Reference
		dataset.Reference = in.Reference
	}

	if len(in.Creators) > 0 {
		// TODO: management many Creators
		dataset.Creators = in.Creators
	}

	if err := s.storage.Store(&dataset,
		option.UpdateWith("UserID", uint(curUsr.ID)),
	); err != nil {
		return nil, err
	}

	return &dataset, nil
}
