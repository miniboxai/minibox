package server

import (
	"io"

	"minibox.ai/minibox/pkg/api/v1/proto"
	"minibox.ai/minibox/pkg/api/v1/types"
)

type StreamWriter struct {
	Filename string
	first    bool
	stream   proto.DatasetService_PutDatasetObjectClient
}

type StreamReader struct {
	buf    []byte
	off    int
	thr    int
	stream proto.DatasetService_GetDatasetObjectClient
}

func NewStreamWriter(filename string, stream proto.DatasetService_PutDatasetObjectClient) io.Writer {
	return &StreamWriter{Filename: filename, stream: stream}
}

func (s *StreamWriter) firstReq() *types.PutObjectRequest {
	var (
		ns, blob = path2ns(s.Filename)
		id       = types.ObjectID{Value: &types.ObjectID_Name{Name: ns}}
		key      = types.PutObjectRequest_Key{Id: &id, Blob: blob}
		value    = types.PutObjectRequest_Key_{Key: &key}
		req      = types.PutObjectRequest{Value: &value}
	)

	return &req
}

func (s *StreamWriter) chunkReq(buf []byte, pos int) *types.PutObjectRequest {
	var (
		chunk = types.PutObjectRequest_Chunk{Data: buf, Position: int64(pos)}
		value = types.PutObjectRequest_Chunk_{Chunk: &chunk}
		req   = types.PutObjectRequest{Value: &value}
	)

	return &req
}

func (s *StreamWriter) Write(buf []byte) (int, error) {
	var req *types.PutObjectRequest
	if !s.first {
		req = s.firstReq()
		s.first = true
		if err := s.stream.Send(req); err != nil {
			return 0, err
		}
	}

	req = s.chunkReq(buf, -1)
	if err := s.stream.Send(req); err != nil {
		return 0, err
	}

	return len(buf), nil
}

func (s *StreamWriter) Close() error {
	// s.stream.CloseAndRecv()
	return nil
}

func NewStreamReader(stream proto.DatasetService_GetDatasetObjectClient) io.Reader {
	return &StreamReader{stream: stream, thr: 1024}
}

func (s *StreamReader) pullData() error {
	reply, err := s.stream.Recv()
	if reply != nil && len(reply.Data) > 0 {
		s.buf = s.buf[s.off:]
		s.buf = append(s.buf, reply.Data...)
		s.off = 0
	}

	return err
}

func (s *StreamReader) Read(p []byte) (n int, err error) {
	l := cap(p)
	b := len(s.buf)

	if b-s.off < s.thr {
		err = s.pullData()
	}

	if l >= b-s.off {
		n = copy(p, s.buf[s.off:])
	} else {
		n = copy(p, s.buf[s.off:s.off+l])
	}
	s.off += n
	return n, err
}
