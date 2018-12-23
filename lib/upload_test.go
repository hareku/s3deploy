package s3deploy

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
)

func TestGetUploadFilePaths(t *testing.T) {
	wd, err := os.Getwd()

	if err != nil {
		t.Fatalf("%+v", errors.Wrap(err, "os.Getwd() failed"))
	}

	uploadPath := filepath.Join(wd, "testdata")

	if err != nil {
		t.Fatalf("%+v", errors.Wrap(err, "filepath.Rel() failed"))
	}

	filePaths, err := GetUploadFilePaths(uploadPath)

	if err != nil {
		t.Fatalf("%+v", errors.Wrap(err, "GetUploadFilePaths() failed"))
	}

	expectPathsCount := 3

	if len(filePaths) != expectPathsCount {
		t.Fatalf(fmt.Sprintf("Get upload file paths does not contain %d files. got %d", expectPathsCount, len(filePaths)))
	}

	expectFiles := []string{
		"index.html",
		"static/main.css",
		"static/main.js",
	}

	for index, expectFile := range expectFiles {
		expect := filepath.Join(wd, "testdata", filepath.FromSlash(expectFile))
		if filePaths[index] != expect {
			t.Errorf("file path does not match expected %q. got=%q", expect, filePaths[index])
		}
	}
}
