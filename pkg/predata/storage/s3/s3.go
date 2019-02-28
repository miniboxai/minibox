package s3

import "github.com/qor/oss/s3"

type (
	Config = s3.Config
	Client = s3.Client
)

func New(cfg Config) *Client {
	return s3.New(cfg)
}
