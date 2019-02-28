package server

import (
	"io"
	"path"

	"minibox.ai/pkg/api/v1/proto"
	"minibox.ai/pkg/api/v1/types"
	bytes "qiniupkg.com/x/bytes.v7"
)

type WriteBuffer = bytes.Buffer

var defaultPartSize = 10 * 1024 * 1024

// GetDatasetObject 获取数据集的数据
// Example:
//     client.GetDatasetObject()
func (s *Server) GetDatasetObject(in *types.GetObjectRequest, stream proto.DatasetService_GetDatasetObjectServer) error {
	var (
		ds        types.Dataset
		partSize  int
		err       error
		datasetId uint32
	)

	datasetId = in.DatasetId.GetDatasetId()
	if datasetId > 0 {
		ds.ID = datasetId
	} else {
		ds.Name = in.DatasetId.GetName()
	}

	if err := s.storage.Load(&ds); err != nil {
		return err
	}

	if in.PartSize > 0 {
		partSize = int(in.PartSize)
	} else {
		partSize = defaultPartSize
	}

	ns := path.Join(ds.Namespace, in.Blob)

	wb := bytes.NewBuffer()
	_, err = s.objects.GetObjectManager(ns, wb)

	var (
		buf = make([]byte, partSize)
		pos int64
		n   int
	)

	for {
		var or types.GetObjectReply
		n, err = wb.ReadAt(buf, pos)
		if n == 0 {
			return io.EOF
		}

		pos += int64(n)

		or.Data = buf[:n]
		or.Position = pos

		stream.Send(&or)
		if err == io.EOF {
			break
		}
	}
	return nil
}
