package server

import (
	"errors"
	"io"
	"log"
	"path"
	"strings"

	"minibox.ai/pkg/api/v1/proto"
	"minibox.ai/pkg/api/v1/types"
	"minibox.ai/pkg/utils"
)

func path2ns(pat string) (ns, dir string) {
	var ss = strings.SplitN(pat, "/", 2)
	if len(ss) == 2 {
		return ss[0], ss[1]
	} else {
		return "__public__", ss[0]
	}
}

func (s *Server) getObject(id *types.ObjectID) (*types.Dataset, error) {
	var (
		datasetId   = id.GetDatasetId()
		datasetName = id.GetName()
		ds          types.Dataset
	)

	if datasetId > 0 {
		ds.ID = datasetId
	} else {
		ns, dir := path2ns(datasetName)
		ds.Name = dir
		ds.Namespace = ns
	}

	if err := s.storage.Load(&ds); err != nil {
		return nil, err
	}

	return &ds, nil
}

func (s *Server) PutDatasetObject(stream proto.DatasetService_PutDatasetObjectServer) error {
	var (
		req   *types.PutObjectRequest
		err   error
		ds    *types.Dataset
		reply types.PutObjectReply
		blob  string
	)

	if req, err = stream.Recv(); err != nil {
		return err
	}

	if key := req.GetKey(); key == nil {
		return errors.New("first packet must Key")
	} else {
		var (
			id = key.GetId()
		)
		blob = key.GetBlob()

		if ds, err = s.getObject(id); err != nil {
			return err
		}
	}

	ns := path.Join(ds.Namespace, blob)

	r, w := io.Pipe()

	var ch = make(chan error)
	go func() {
		_, err := s.objects.PutObjectManager(ns, r)
		ch <- err
	}()

	var (
		size int
		exit bool
	)

	for !exit {
		if req, err = stream.Recv(); err == io.EOF {
			w.Close()
			break
		}

		chunk := req.GetChunk()
		log.Printf("chunk: %s", utils.Prettify(chunk.Data))

		if chunk != nil {
			if n, err := w.Write(chunk.Data); err != nil {
				return err
			} else {
				size += n
			}
		}
	}
	if err = <-ch; err != nil {
		return err
	}

	reply.Size_ = uint32(size)
	// b.Close()
	return stream.SendAndClose(&reply)
}
