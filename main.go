package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	s3deploy "github.com/hareku/s3deploy/lib"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "s3-safely-deploy"
	app.Usage = "This tool safely deploy to Amazon S3"
	app.Version = "0.2.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "bucket",
			Usage: "deploy bucket name",
		},
		cli.StringFlag{
			Name:  "directory",
			Value: "/",
			Usage: "deploy bucket directory",
		},
		cli.StringFlag{
			Name:  "region",
			Value: "ap-northeast-1",
			Usage: "aws region",
		},
		cli.StringFlag{
			Name:  "versions-file",
			Value: "__s3_safety_deploy_versions",
			Usage: "versions manegement file name",
		},
	}

	app.Action = func(context *cli.Context) error {
		if context.String("bucket") == "" {
			return errors.New("bucket option is required")
		}

		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(context.String("region")),
		}))

		config := &s3deploy.Config{
			S3Client:         s3.New(sess),
			Bucket:           context.String("bucket"),
			Prefix:           context.String("directory"),
			VersionsFileName: "__s3_safety_deploy_versions.json",
			UploadPath:       context.Args().First(),
		}

		return s3deploy.Deploy(config)
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatalf("%v", errors.Wrap(err, "s3-safely-deploy error"))
	}
}
