# S3 Safely Deploy (Golang)

## Description
A simple tool to safely deploy to Amazon S3.

This tool manages three versions of your application.

## How to deploy
S3 has 3 versions objects.

- Current Version (Deploying version)
- Previous Versoin (The second last version)
- NextDelete Version (The third last version)

1. Delete files of NextDelete that are not in Current and Previous
4. Deploy files.
2. Set Previous as NextDelete and Current as Previous
3. Set deploying files to Current, and put versions management file (__s3_safety_deploy_versions.json) to S3

## Usage
Install
`go get -u github.com/hareku/s3deploy`

Deploy
`s3deploy --bucket my-bucket-name ./path-to-your-app`

Options
```
GLOBAL OPTIONS:
   --bucket value         deploy bucket name
   --directory value      deploy bucket directory (default: "/")
   --region value         aws region (default: "ap-northeast-1")
```

## Example of use
I developed this tool as it became necessary for Nuxt.js, Ruby on Rails projects.

### Nuxt.js Static Website
[How to safely deploy static website to Amazon S3 (Qiita)](https://qiita.com/hareku/items/acd9dfe5d77293a4a6d9#%E5%AE%89%E5%85%A8%E3%81%AA%E3%83%87%E3%83%97%E3%83%AD%E3%82%A4%E3%82%92%E3%81%A9%E3%81%AE%E3%82%88%E3%81%86%E3%81%AB%E3%81%99%E3%82%8B%E3%81%8B)

### ECS Application
[Avoid the problem that the static file becomes 404 with constant probability when ECS is deployed (Qiita)](https://qiita.com/hareku/items/6be1b71e58033b9739fd)

## Contribution
1. Fork this project
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request
