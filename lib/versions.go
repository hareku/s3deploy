package s3deploy

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

type Versions struct {
	Current    *Version `json:"Current"`
	Previous   *Version `json:"Previous"`
	NextDelete *Version `json:"NextDelete"`
}

type Version struct {
	Keys []string `json:"Keys"`
}

func (config *Config) GetVersions() (*Versions, error) {
	versions := &Versions{
		Current:    &Version{Keys: []string{}},
		Previous:   &Version{Keys: []string{}},
		NextDelete: &Version{Keys: []string{}},
	}

	s3Result, err := config.S3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(config.Bucket),
		Key:    aws.String(strings.TrimRight(config.Prefix, "/") + "/" + config.VersionsFileName),
	})

	if err != nil {
		return versions, nil
	}

	defer s3Result.Body.Close()

	byteArray, err := ioutil.ReadAll(s3Result.Body)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to read a body of S3 versions file")
	}

	if err := json.Unmarshal(byteArray, versions); err != nil {
		return versions, errors.Wrap(err, "Failed to unmarshal versions file json")
	}

	return versions, err
}

func (config *Config) SaveVersions(versions *Versions) error {
	jsonBytes, err := json.Marshal(versions)

	if err != nil {
		return errors.Wrap(err, "Failed to marshal versions file json")
	}

	_, err = config.S3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(config.Bucket),
		Key:    aws.String(strings.TrimRight(config.Prefix, "/") + "/" + config.VersionsFileName),
		Body:   bytes.NewReader(jsonBytes),
	})

	if err != nil {
		return errors.Wrap(err, "Failed to put versions file to S3")
	}

	return nil
}

func (versions *Versions) MigrateVersion(uploadVersion *Version) {
	versions.NextDelete = versions.Previous
	versions.Previous = versions.Current
	versions.Current = uploadVersion
}

func (versions *Versions) PrintVersions() {
	versions.Current.printVersion("CurrentVersion")
	versions.Previous.printVersion("PreviousVersion")
	versions.NextDelete.printVersion("NextDeleteVersion")
}

func (version *Version) printVersion(versionName string) {
	log.Printf("== %s ==\n", versionName)

	if len(version.Keys) == 0 {
		log.Printf("%s is empty\n", versionName)
		return
	}

	for _, key := range version.Keys {
		log.Print(key)
	}
}
