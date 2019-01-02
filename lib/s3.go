package s3deploy

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (config *Config) ResolveDeleteObjects(versions *Versions) []*s3.ObjectIdentifier {
	var objects []*s3.ObjectIdentifier

	for _, nextDeleteVersionKey := range versions.NextDelete.Keys {
		if versionContainKey(versions.Current, nextDeleteVersionKey) || versionContainKey(versions.Previous, nextDeleteVersionKey) {
			continue
		}

		if strings.HasSuffix(nextDeleteVersionKey, config.VersionsFileName) {
			continue
		}

		objects = append(objects, &s3.ObjectIdentifier{
			Key: &nextDeleteVersionKey,
		})
	}

	return objects
}

func versionContainKey(version *Version, key string) bool {
	for _, versionKey := range version.Keys {
		if versionKey == key {
			return true
		}
	}

	return false
}

func (config *Config) DeleteObjects(objects []*s3.ObjectIdentifier) (*s3.DeleteObjectsOutput, error) {
	for _, object := range objects {
		log.Printf("Deleting: %s", *object.Key)
	}

	return config.S3Client.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(config.Bucket),
		Delete: &s3.Delete{
			Objects: objects,
		},
	})
}
