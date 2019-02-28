package server

import (
	"context"

	"minibox.ai/pkg/api/v1/types"
	"minibox.ai/pkg/logger"
	"minibox.ai/pkg/server/internal/acl"
	"minibox.ai/pkg/server/internal/option"
)

func (s *Server) ListDatasets(ctx context.Context, in *types.ListDatasetsRequest) (*types.DatasetsReply, error) {
	var (
		datasets = make([]*types.Dataset, 0, 30)
		filters  = option.
				NewFilters().
				Field("namespace", in.Namespace)
		ok     bool
		curUsr *types.User
		role   *Role
	)

	logger.S().Infow("ListDatasets",
		"ns", in.Namespace)
	curUsr, ok = s.currentUser(ctx)
	if !ok {
		return nil, ErrInvalidUser
	}

	role = &Role{usr: curUsr}

	if !role.Can(acl.ListDatasets) {
		return nil, ErrInvalidPermission(acl.ListDatasets)
	} else if in.ShowPrivate && (curUsr.Namespace == in.Namespace || role.Can(acl.ListPrivateDatasets)) {
		// filters.Contains("private", []bool{true, false})
	} else if in.ShowPrivate {
		return nil, ErrInvalidPermission(acl.ListPrivateDatasets)
	} else {
		// filters.NotField("private", true).OrField("private", nil)
	}

	if err := s.storage.List(&datasets,
		option.WithSub("User"),
		option.Filters(filters),
	); err != nil {
		return nil, err
	}

	var reply = types.DatasetsReply{Datasets: datasets}

	return &reply, nil
}
