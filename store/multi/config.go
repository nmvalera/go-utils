package multistore

import (
	filestore "github.com/kkrt-labs/go-utils/store/file"
	s3store "github.com/kkrt-labs/go-utils/store/s3"
)

type Config struct {
	FileConfig *filestore.Config
	S3Config   *s3store.Config
}
