package main

import (
  "fmt"
  "github.com/gographics/imagick/imagick"
  "github.com/golang/groupcache/lru"
  "launchpad.net/goamz/aws"
  "launchpad.net/goamz/s3"
  "net"
  "net/http"
  "net/http/fcgi"
  "strconv"
  "time"
)

var cache *lru.Cache
var bucket *s3.Bucket

func sec() int64 {
  return time.Now().UnixNano()
}

func GetImage(path string) (bytes []byte, err error) {
  blob, res := cache.Get(path)
  if !res {
    s3blob, err := bucket.Get(path)
    if err != nil {
      return nil, err
    }
    cache.Add(path, s3blob)

    return s3blob, nil
  }
  return blob.([]byte), nil
}

func gofuHandler(writer http.ResponseWriter, req *http.Request) {
  t := sec()
  imageBlob, err := GetImage(req.URL.Path[1:])

  if err != nil {
    fmt.Println(err)
    http.NotFound(writer, req)
    return
  }
  fmt.Printf("Get by S3:        %d\n", sec()-t)

  t = sec()
  magick_wand := imagick.NewMagickWand()
  defer func() {
    magick_wand.Destroy()
  }()
  magick_wand.ReadImageBlob(imageBlob)
  fmt.Printf("MagickWand:       %d\n", sec()-t)

  t = sec()
  query := req.URL.Query()

  width := magick_wand.GetImageWidth()
  height := magick_wand.GetImageHeight()

  if query["w"] != nil {
    w, _ := strconv.Atoi(query["w"][0])
    width = uint(w)
  }
  if query["h"] != nil {
    h, _ := strconv.Atoi(query["h"][0])
    height = uint(h)
  }
  magick_wand.ResizeImage(width, height, imagick.FILTER_CUBIC, 1)
  fmt.Printf("Parse and Resize: %d\n", sec()-t)

  t = sec()
  responseBlob := magick_wand.GetImageBlob()

  writer.Header().Add("Content-Type", http.DetectContentType(responseBlob))
  writer.Header().Add("Content-Length", strconv.Itoa(len(responseBlob)))
  writer.Write(responseBlob)
  fmt.Printf("Send Response:    %d\n", sec()-t)
}

func startWithHttp() {
  address := fmt.Sprintf("%s:%d", gofu_config.Bind, gofu_config.Port)

  http.HandleFunc("/", gofuHandler)
  http.ListenAndServe(address, nil)
}

func startWithFcgi() {
  address := fmt.Sprintf("%s:%d", gofu_config.Bind, gofu_config.Port)

  mux := http.NewServeMux()
  mux.HandleFunc("/", gofuHandler)
  listen, _ := net.Listen("tcp", address)
  fcgi.Serve(listen, mux)
}

func start() {
  cache = lru.New(gofu_config.MaxCache)
  s3client := s3.New(gofu_config.S3, aws.APNortheast)
  bucket = s3client.Bucket(gofu_config.Bucket)

  if gofu_config.Fcgi {
    startWithFcgi()
  } else {
    startWithHttp()
  }
}
