package gofu

import (
  "github.com/gographics/imagick/imagick"
  "launchpad.net/goamz/aws"
  "launchpad.net/goamz/s3"
  "runtime"
)

type Image struct {
  Path      string
  Processed bool
  Blob      []byte
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
    Path:      path,
    Processed: false,
  }, err
}

func getByS3(path string) ([]byte, error) {
  client := <-s3clients
  defer func() { s3clients <- client }()

  return client.Get(path)
}

func (i *Image) Process(wand *imagick.MagickWand, q *Query) {
  defer wand.Clear()

  // wand.SetOption("jpeg:size", fmt.Sprintf("%dx%d", q.ResizedWidth, q.ResizedHeight))

  bench("wand read", func() { wand.ReadImage(i.Path) })
  bench("wand set image format", func() { wand.SetImageFormat("jpeg") })
  bench("wand set compress", func() { wand.SetCompression(imagick.COMPRESSION_JPEG2000) })
  bench("wand set image quality", func() { wand.SetImageCompressionQuality(95) })

  bench("wand resize", func() {
    i.resize(wand, q.ResizedWidth, q.ResizedHeight)
  })

  bench("wand strip", func() {
    wand.StripImage()
  })

  i.Blob = wand.GetImageBlob()
  i.Processed = true
}

func (i *Image) resize(wand *imagick.MagickWand, w, h uint) {
  ow := wand.GetImageWidth()
  oh := wand.GetImageHeight()

  if (float64(ow) / float64(oh)) < (float64(w) / float64(h)) {
    h = oh * w / ow
  } else {
    w = ow * h / oh
  }

  wand.SetImageInterpolateMethod(imagick.INTERPOLATE_PIXEL_BICUBIC)
  wand.ResizeImage(w, h, imagick.FILTER_LANCZOS2_SHARP, 1)
}
