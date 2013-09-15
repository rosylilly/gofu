package gofu

import "launchpad.net/goamz/aws"

type S3Config struct {
  BucketName string `json:"bucket_name"`
  Region     string `json:"region"`
  EndPoint   string `json:"end_point"`
  AccessKey  string `json:"access_key"`
  SecretKey  string `json:"secret_key"`
}

func (c *S3Config) AwsRegion() aws.Region {
  region, exist := aws.Regions[c.Region]
  if !exist {
    region = aws.Region{
      Name:                 "UserDefined",
      EC2Endpoint:          "",
      S3Endpoint:           c.EndPoint,
      S3BucketEndpoint:     c.EndPoint,
      S3LocationConstraint: false,
      S3LowercaseBucket:    false,
      SDBEndpoint:          "",
      SNSEndpoint:          "",
      SQSEndpoint:          "",
      IAMEndpoint:          "",
    }
  }

  return region
}
