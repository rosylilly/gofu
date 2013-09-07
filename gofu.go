package main

import (
  "github.com/gographics/imagick/imagick"
  "launchpad.net/goamz/aws"
  "launchpad.net/goamz/s3"
  "net"
  "net/http"
  "net/http/fcgi"
  "strconv"
  "fmt"
)

var bucket *s3.Bucket

func gofuHandler(writer http.ResponseWriter, req *http.Request) {
  imageBlob, err := bucket.Get(req.URL.Path[1:])

  if(err != nil) {
    fmt.Println(err)
    http.NotFound(writer, req)
    return
  }

  magick_wand := imagick.NewMagickWand()
  defer func(){
    magick_wand.Destroy()
  }()
  magick_wand.ReadImageBlob(imageBlob)

  query := req.URL.Query()

  width := magick_wand.GetImageWidth()
  height := magick_wand.GetImageHeight()

  if(query["w"] != nil) {
    w, _ := strconv.Atoi(query["w"][0])
    width = uint(w)
  }
  if(query["h"] != nil) {
    h, _ := strconv.Atoi(query["h"][0])
    height = uint(h)
  }
  magick_wand.ResizeImage(width, height, imagick.FILTER_CUBIC, 1)


  responseBlob := magick_wand.GetImageBlob()

  writer.Header().Add("Content-Type", http.DetectContentType(responseBlob))
  writer.Header().Add("Content-Length", strconv.Itoa(len(responseBlob)))
  writer.Write(responseBlob)
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
  s3client := s3.New(gofu_config.S3, aws.APNortheast)
  bucket = s3client.Bucket(gofu_config.Bucket)

  if(gofu_config.Fcgi) {
    startWithFcgi()
  } else {
    startWithHttp()
  }
}
