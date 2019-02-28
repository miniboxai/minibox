package server

import (
	"context"
	"path"

	"minibox.ai/minibox/pkg/api/v1/types"
)

func (s *Server) ListDatasetObjects(ctx context.Context, in *types.ListObjectsRequest) (*types.ObjectsReply, error) {
	var (
		curUsr      *types.User
		ok          bool
		ns          string
		reply       types.ObjectsReply
		datasetId   uint32 = in.Id.GetDatasetId()
		datasetName string = in.Id.GetName()
		dataset     types.Dataset
	)

	curUsr, ok = s.currentUserWith(ctx, LoadOrgs())
	if !ok {
		return nil, ErrInvalidUser
	}

	if datasetId > 0 {
		dataset = types.Dataset{
			ID: datasetId,
		}

		// load dataset
		if err := s.storage.Load(&dataset); err != nil {
			return nil, err
		}
	} else {
		dataset = types.Dataset{
			Name: datasetName,
		}

		// load dataset
		if err := s.storage.Load(&dataset); err != nil {
			return nil, err
		}
	}

	if dataset.Namespace == curUsr.Namespace {
		goto Exec
	} else if nss := s.orgsNamespaces(curUsr.Organizations); !s.inNamespaces(nss, dataset.Namespace) {
		return nil, ErrInvalidNamespace(dataset.Namespace)
	}

Exec:
	ns = path.Join(dataset.Namespace, dataset.Name)
	objs, err := s.objects.ListObjects(ns)
	if err != nil {
		return nil, err
	}

	reply.Objects = make([]*types.Object, len(objs))
	for i, o := range objs {
		reply.Objects[i] = &types.Object{
			Etag:      o.ETag,
			Name:      o.Key,
			Size_:     uint64(o.Size),
			UpdatedAt: Time(o.LastModified),
		}
	}
	return &reply, nil
}
