package s3deploy

import "github.com/aws/aws-sdk-go/service/s3"

type Config struct {
	S3Client         *s3.S3
	Bucket           string
	Prefix           string
	VersionsFileName string
	UploadPath       string
}
