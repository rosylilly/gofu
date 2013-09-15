package gofu

import (
  "launchpad.net/goamz/aws"
  "launchpad.net/goamz/s3"
  "runtime"
)

type Image struct {
  Path string
  Blob []byte
}

var s3clients chan *s3.Bucket

func InitS3Client(c S3Config) {
  auth := aws.Auth{c.AccessKey, c.SecretKey}
  s3client := s3.New(auth, c.AwsRegion())

  s3clients = make(chan *s3.Bucket, runtime.NumCPU())
  for i := 0; i < runtime.NumCPU(); i++ {
    bucket := s3client.Bucket(c.BucketName)
    s3clients <- bucket
  }
}

func GetImage(path string) (*Image, error) {
  path, err := cache.Fetch(path, getByS3)

  return &Image{
    Path: path,
  }, err
}

func getByS3(path string) ([]byte, error) {
  client := <-s3clients
  defer func() { s3clients <- client }()

  return client.Get(path)
}
