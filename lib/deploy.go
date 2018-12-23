package s3deploy

import (
	"github.com/pkg/errors"
)

func Deploy(config *Config) error {
	versions, err := config.GetVersions()
	if err != nil {
		return errors.Wrap(err, "Failed to get versions file from S3")
	}

	bucketObjects, err := config.ListBucketObjects()
	if err != nil {
		return errors.Wrap(err, "Failed to list S3 objects")
	}

	deleteObjects := config.ResolveDeleteObjects(bucketObjects, versions)
	if len(deleteObjects) > 0 {
		_, err = config.DeleteObjects(deleteObjects)
		if err != nil {
			return errors.Wrap(err, "Failed to delete old version objects from S3")
		}
	}

	uploadFilePaths, err := GetUploadFilePaths(config.UploadPath)
	if err != nil {
		return errors.Wrap(err, "Failed to list upload file paths")
	}

	newVersion, err := config.UploadFiles(uploadFilePaths)
	if err != nil {
		return errors.Wrap(err, "Failed to upload new files to S3")
	}

	versions.MigrateVersion(newVersion)
	err = config.SaveVersions(versions)
	if err != nil {
		return errors.Wrap(err, "Failed to save versions file to S3")
	}

	return nil
}
