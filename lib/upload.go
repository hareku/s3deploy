package s3deploy

import (
	"bytes"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetUploadFilePaths(uploadPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(uploadPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "Prevent panic by handling failure accessing a path %q", path)
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return files, errors.Wrapf(err, "error walking the path %q", uploadPath)
	}

	return files, nil
}

func (config *Config) UploadFiles(uploadFilePaths []string) (*Version, error) {
	version := &Version{
		Keys: make([]string, len(uploadFilePaths)),
	}

	for index, filePath := range uploadFilePaths {
		fileBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return version, errors.Wrapf(err, "Failed to read upload file: %s", filePath)
		}

		relPath, err := filepath.Rel(config.UploadPath, filePath)

		if err != nil {
			return version, errors.Wrapf(err, "Failed to resolve upload file rel path: %s", filePath)
		}

		objectKey := strings.TrimLeft(strings.TrimRight(config.Prefix, "/")+"/"+filepath.ToSlash(relPath), "/")

		contentType := mime.TypeByExtension(filepath.Ext(relPath))
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		_, err = config.S3Client.PutObject(&s3.PutObjectInput{
			Bucket:      aws.String(config.Bucket),
			Key:         aws.String(objectKey),
			Body:        bytes.NewReader(fileBytes),
			ContentType: &contentType,
		})

		log.Printf("Uploaded: %s", objectKey)

		if err != nil {
			return version, errors.Wrap(err, "Failed to put versions file to s3")
		}

		version.Keys[index] = objectKey
	}

	return version, nil
}
