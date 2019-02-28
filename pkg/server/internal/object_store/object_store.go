package object_store

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"minibox.ai/minibox/pkg/logger"
	"minibox.ai/minibox/pkg/utils"
)

type WriteAtBuffer = aws.WriteAtBuffer

type ObjectStore struct {
	APPID       string
	SECRET      string
	Region      string
	Endpoint    string
	StoreBucket string
	Root        string
}

type Bucket struct {
	Name      string
	CreatedAt time.Time
}

type Object struct {
	ETag         string    `type:"string"`
	Key          string    `min:"1" type:"string"`
	LastModified time.Time `type:"timestamp"`
	Size         int64     `type:"integer"`
	// Owner        *Owner     `type:"structure"`
	StorageClass string `type:"string" enum:"ObjectStorageClass"`
}

func NewObjectStore(appid, secret string, opts ...OptionFunc) *ObjectStore {
	objs := &ObjectStore{APPID: appid, SECRET: secret}
	for _, op := range opts {
		op(objs)
	}
	return objs
}

func (objs *ObjectStore) config() *aws.Config {
	cfg := aws.NewConfig()
	cfg.WithRegion(objs.Region)
	cfg.WithMaxRetries(3)
	cfg.WithEndpoint(objs.Endpoint)
	cfg.WithCredentials(
		credentials.NewStaticCredentials(objs.APPID, objs.SECRET, ""))

	return cfg
}

func (objs *ObjectStore) ListBuckets() ([]*Bucket, error) {
	sess := session.Must(session.NewSession(objs.config()))
	svc := s3.New(sess)

	logger.S().Infow("ListBuckets",
		"region", objs.Region,
		"bucket", objs.StoreBucket)
	output, err := svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}
	var results = make([]*Bucket, len(output.Buckets))

	for i, bucket := range output.Buckets {
		results[i] = &Bucket{
			CreatedAt: *bucket.CreationDate,
			Name:      *bucket.Name,
		}
	}

	return results, nil
}

func (objs *ObjectStore) ListObjects(pat string) ([]*Object, error) {
	var (
		sess = session.Must(session.NewSession(objs.config()))
		svc  = s3.New(sess)
		ctx  = aws.BackgroundContext()
		in   = s3.ListObjectsV2Input{
			// in := s3.ListObjectsInput{
			Bucket: aws.String(objs.StoreBucket),
			// StartAfter: aws.String(pat),
			Prefix: aws.String(path.Join(objs.Root, pat)),
			// Delimiter: aws.String(""),
			MaxKeys: aws.Int64(1024),
		}
		results []*Object
	)

	logger.S().Info("ListObjects",
		"region", objs.Region,
		"bucket", in.Bucket,
		"max_keys", in.MaxKeys,
		"prefix", in.Prefix,
	)

	if err := svc.ListObjectsV2PagesWithContext(ctx, &in, func(output *s3.ListObjectsV2Output, last bool) bool {
		results = make([]*Object, len(output.Contents))

		for i, obj := range output.Contents {
			results[i] = &Object{
				ETag:         *obj.ETag,
				Key:          *obj.Key,
				LastModified: *obj.LastModified,
				Size:         *obj.Size,
				StorageClass: *obj.StorageClass,
			}
		}

		return false
	}); err != nil {
		return nil, err
	}

	return results, nil
}

func (objs *ObjectStore) DeleteObject(pat string, opts ...ObjectFunc) (bool, error) {
	var (
		sess = session.Must(session.NewSession(objs.config()))
		svc  = s3.New(sess)
		in   = s3.DeleteObjectInput{
			Bucket: aws.String(objs.StoreBucket),
			Key:    aws.String(path.Join(objs.Root, pat)),
		}
		opt = new(ObjectOption)
	)

	for _, op := range opts {
		op(opt)
	}

	if len(opt.Version) > 0 {
		in.VersionId = aws.String(opt.Version)
	}

	output, err := svc.DeleteObject(&in)
	var marker bool
	if output.DeleteMarker != nil {
		marker = *output.DeleteMarker
	}
	return marker, err
}

func (objs *ObjectStore) setRange(offset, size int) *string {
	var rng string

	if offset == 0 && size == 0 {
		return nil
	}

	if offset > 0 && size == 0 {
		rng = fmt.Sprintf("bytes %d-", offset)
		return &rng
	} else if size > 0 {
		rng = fmt.Sprintf("bytes %d-%d", offset, size)
		return &rng
	} else {
		return nil
	}
}

func (objs *ObjectStore) GetObject(pat string, offset, size int, opts ...ObjectFunc) (io.ReadCloser, error) {
	var (
		sess = session.Must(session.NewSession(objs.config()))
		svc  = s3.New(sess)
		in   = s3.GetObjectInput{
			Bucket: aws.String(objs.StoreBucket),
			Key:    aws.String(path.Join(objs.Root, pat)),
		}
		opt = new(ObjectOption)
	)

	for _, op := range opts {
		op(opt)
	}

	if len(opt.Version) > 0 {
		in.VersionId = aws.String(opt.Version)
	}

	in.Range = objs.setRange(offset, size)

	output, err := svc.GetObject(&in)
	if err != nil {
		return nil, err
	}
	return output.Body, nil
}

func (objs *ObjectStore) GetObjectManager(pat string, w io.WriterAt, opts ...ObjectFunc) (int, error) {
	var (
		sess       = session.Must(session.NewSession(objs.config()))
		downloader = s3manager.NewDownloader(sess)
		in         = s3.GetObjectInput{
			Bucket: aws.String(objs.StoreBucket),
			Key:    aws.String(path.Join(objs.Root, pat)),
		}
		opt = new(ObjectOption)
	)

	for _, op := range opts {
		op(opt)
	}
	log.Printf("download: %s", utils.Prettify(&in))
	writen, err := downloader.Download(w, &in)
	if err != nil {
		return 0, err
	}
	return int(writen), nil
}

func (objs *ObjectStore) PutObjectManager(pat string, rd io.Reader, opts ...ObjectFunc) (string, error) {
	var (
		sess     = session.Must(session.NewSession(objs.config()))
		uploader = s3manager.NewUploader(sess)
		in       = s3manager.UploadInput{
			Bucket: aws.String(objs.StoreBucket),
			Key:    aws.String(path.Join(objs.Root, pat)),
			Body:   rd,
		}

		opt = new(ObjectOption)
	)

	for _, op := range opts {
		op(opt)
	}

	if len(opt.Tag) > 0 {
		in.Tagging = aws.String(opt.Tag)
	}

	log.Printf("Key: %s", *in.Key)
	// in.Range = objs.setRange(offset, size)

	output, err := uploader.Upload(&in)
	// output, err := svc.PutObject(&in)
	if err != nil {
		return "", err
	}

	var versionId string

	if output.VersionID != nil {
		versionId = *output.VersionID
	}

	return versionId, nil
}

func (objs *ObjectStore) CopyObject(src, dst string, opts ...ObjectFunc) error {
	var (
		sess = session.Must(session.NewSession(objs.config()))
		svc  = s3.New(sess)
		in   = s3.CopyObjectInput{
			Bucket:     aws.String(objs.StoreBucket),
			CopySource: aws.String(path.Join(objs.StoreBucket, objs.Root, src)),
			Key:        aws.String(path.Join(objs.Root, dst)),
		}
		opt = new(ObjectOption)
	)

	for _, op := range opts {
		op(opt)
	}

	if len(opt.Tag) > 0 {
		in.Tagging = aws.String(opt.Tag)
	}

	_, err := svc.CopyObject(&in)
	return err
}

func (objs *ObjectStore) MoveObject(src, dst string, opts ...ObjectFunc) error {
	var (
		sess = session.Must(session.NewSession(objs.config()))
		svc  = s3.New(sess)
		in   = s3.CopyObjectInput{
			Bucket:     aws.String(objs.StoreBucket),
			CopySource: aws.String(path.Join(objs.StoreBucket, objs.Root, src)),
			Key:        aws.String(path.Join(objs.Root, dst)),
		}
		opt = new(ObjectOption)
	)

	for _, op := range opts {
		op(opt)
	}

	if len(opt.Tag) > 0 {
		in.Tagging = aws.String(opt.Tag)
	}

	if _, err := svc.CopyObject(&in); err != nil {
		return err
	}

	var (
		delIn = s3.DeleteObjectInput{
			Bucket: aws.String(objs.StoreBucket),
			Key:    aws.String(path.Join(objs.Root, src)),
		}
	)

	if _, err := svc.DeleteObject(&delIn); err != nil {
		return err
	}

	return nil
}

func (objs *ObjectStore) GetObjectPresignURL(pat string, opts ...ObjectFunc) (url string, signedHeaders http.Header, err error) {
	var (
		sess = session.Must(session.NewSession(objs.config()))
		svc  = s3.New(sess)
		in   = s3.GetObjectInput{
			Bucket: aws.String(objs.StoreBucket),
			Key:    aws.String(path.Join(objs.Root, pat)),
		}
		opt    = new(ObjectOption)
		expire = 15 * time.Minute
		req    *request.Request
	)

	for _, op := range opts {
		op(opt)
	}

	req, _ = svc.GetObjectRequest(&in)

	if opt.Expire > 0 {
		expire = opt.Expire
	}
	url, signedHeaders, err = req.PresignRequest(expire)
	return
}

func (objs *ObjectStore) PutObjectPresignURL(pat string, opts ...ObjectFunc) (url string, signedHeaders http.Header, err error) {
	var (
		sess = session.Must(session.NewSession(objs.config()))
		svc  = s3.New(sess)
		in   = s3.PutObjectInput{
			Bucket: aws.String(objs.StoreBucket),
			Key:    aws.String(path.Join(objs.Root, pat)),
		}
		opt    = new(ObjectOption)
		expire = 15 * time.Minute
		req    *request.Request
	)

	for _, op := range opts {
		op(opt)
	}

	req, _ = svc.PutObjectRequest(&in)

	if opt.Expire > 0 {
		expire = opt.Expire
	}
	url, signedHeaders, err = req.PresignRequest(expire)
	return
}

func (objs *ObjectStore) ExistsObject(pat string) (bool, error) {
	var (
		sess = session.Must(session.NewSession(objs.config()))
		svc  = s3.New(sess)
		in   = s3.HeadObjectInput{
			Bucket: aws.String(objs.StoreBucket),
			Key:    aws.String(path.Join(objs.Root, pat)),
		}
	)

	if err := svc.WaitUntilObjectExists(&in); err != nil {
		return false, err
	}

	return true, nil
}

func (objs *ObjectStore) NotExistsObject(pat string) (bool, error) {
	var (
		sess = session.Must(session.NewSession(objs.config()))
		svc  = s3.New(sess)
		in   = s3.HeadObjectInput{
			Bucket: aws.String(objs.StoreBucket),
			Key:    aws.String(path.Join(objs.Root, pat)),
		}
	)

	if err := svc.WaitUntilObjectNotExists(&in); err != nil {
		return false, err
	}

	return true, nil
}

func (objs *ObjectStore) PutObject(pat string, buf []byte, opts ...ObjectFunc) (int, string, error) {
	var (
		sess = session.Must(session.NewSession(objs.config()))
		svc  = s3.New(sess)
		rd   = bytes.NewReader(buf)
		in   = s3.PutObjectInput{
			Bucket: aws.String(objs.StoreBucket),
			Body:   aws.ReadSeekCloser(rd),
			Key:    aws.String(path.Join(objs.Root, pat)),
		}
		opt = new(ObjectOption)
	)

	for _, op := range opts {
		op(opt)
	}

	if len(opt.Tag) > 0 {
		in.Tagging = aws.String(opt.Tag)
	}

	log.Printf("Key: %s", *in.Key)
	// in.Range = objs.setRange(offset, size)

	output, err := svc.PutObject(&in)
	if err != nil {
		return 0, "", err
	}

	var versionId string

	if output.VersionId != nil {
		versionId = *output.VersionId
	}

	return len(buf), versionId, nil
}
