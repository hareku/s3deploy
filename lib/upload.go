package s3deploy

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type UploadFile struct {
	Path        string
	RelPath     string
	ContentType string
}

func (config *Config) GetUploadFiles(uploadPath string) ([]*UploadFile, error) {
	var files []*UploadFile

	err := filepath.Walk(uploadPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "Failed to access the path %q", path)
		}

		if !info.IsDir() {
			file, err := NewUploadFile(path, config.UploadPath)
			if err != nil {
				return errors.Wrapf(err, "Failed to make UploadFile")
			}
			files = append(files, file)
		}

		return nil
	})

	if err != nil {
		return files, errors.Wrapf(err, "Error walking the path %q", uploadPath)
	}

	return files, nil
}

func NewUploadFile(filePath string, basePath string) (*UploadFile, error) {
	file := &UploadFile{
		Path: filePath,
	}

	relPath, err := filepath.Rel(basePath, filePath)

	if err != nil {
		return file, errors.Wrapf(err, "Failed to resolve upload file rel path: %s", file.Path)
	}

	file.RelPath = relPath

	contentType := mime.TypeByExtension(filepath.Ext(file.RelPath))
	if contentType == "" {
		file.ContentType = "application/octet-stream"
	} else {
		file.ContentType = contentType
	}

	return file, nil
}

func resolveUploadObjectKey(prefix string, uploadFile *UploadFile) string {
	slashPath := filepath.ToSlash(uploadFile.RelPath)
	withPrefixPath := strings.TrimRight(prefix, "/") + "/" + slashPath

	return strings.TrimLeft(withPrefixPath, "/")
}

func (config *Config) UploadFiles(uploadFiles []*UploadFile) (*Version, error) {
	version := &Version{
		Keys: make([]string, len(uploadFiles)),
	}

	eg, ctx := errgroup.WithContext(context.Background())

	for index, uploadFile := range uploadFiles {
		uploadFileForArg := uploadFile
		eg.Go(func() error {
			return config.uploadFile(ctx, uploadFileForArg)
		})

		version.Keys[index] = resolveUploadObjectKey(config.Prefix, uploadFile)
	}

	if err := eg.Wait(); err != nil {
		return version, errors.Wrap(err, "Failed to upload files")
	}

	return version, nil
}

func (config *Config) uploadFile(ctx context.Context, uploadFile *UploadFile) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		// do nothing
	}

	errCh := make(chan error, 1)
	defer close(errCh)

	go func() {
		fileBytes, err := ioutil.ReadFile(uploadFile.Path)
		if err != nil {
			errCh <- errors.Wrapf(err, "Failed to read file for upload: %s", uploadFile.Path)
			return
		}

		_, err = config.S3Client.PutObject(&s3.PutObjectInput{
			Bucket:      aws.String(config.Bucket),
			Key:         aws.String(resolveUploadObjectKey(config.Prefix, uploadFile)),
			Body:        bytes.NewReader(fileBytes),
			ContentType: &uploadFile.ContentType,
		})

		if err != nil {
			errCh <- errors.Wrap(err, "Failed to upload file to S3")
			return
		}

		fmt.Printf("Uploaded: %s\n", uploadFile.Path)
		errCh <- nil
		return
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-ctx.Done():
		<-errCh
		return ctx.Err()
	}

	return nil
}
