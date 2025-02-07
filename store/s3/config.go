package s3store

import aws "github.com/kkrt-labs/go-utils/aws"

type Config struct {
	ProviderConfig *aws.ProviderConfig
	Bucket         string
	KeyPrefix      string
}
