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
  "time"
)

type MagickWandChan chan *imagick.MagickWand

type GofuServer struct {
  bind   string
  bucket *s3.Bucket
  lru    *lru.Cache
  logger *GofuLogger
  wands  MagickWandChan
}

func NewGofuServer() *GofuServer {
  server := new(GofuServer)

  server.bind = fmt.Sprintf("%s:%d", gofuConfig.Bind, gofuConfig.Port)
  server.bucket = s3.New(gofuConfig.S3Config, aws.APNortheast).Bucket(gofuConfig.Bucket)
  server.lru = lru.New(gofuConfig.MaxCache)
  server.logger = NewGofuLogger()
  server.wands = make(MagickWandChan, gofuConfig.MaxProc)

  for i := 0; i < gofuConfig.MaxProc; i++ {
    mw := imagick.NewMagickWand()
    server.wands <- mw
  }

  return server
}

func (server *GofuServer) handlerSelector(req *http.Request) func(res *GofuResponse, req *http.Request) {
  if req.Method == "GET" || req.Method == "HEAD" {
    return server.imageHandler
  }

  return func(res *GofuResponse, req *http.Request) {
    res.Status = http.StatusMethodNotAllowed
  }
}

func (server *GofuServer) mainHandler(w http.ResponseWriter, req *http.Request) {
  requestTime := time.Now()
  res := NewGofuResponse(w)
  defer func() {
    if err := recover(); err != nil {
      res.Status = http.StatusInternalServerError
      res.Body = []byte("")
    }

    res.Write()
    server.logger.log(res, req, requestTime)
  }()

  server.handlerSelector(req)(res, req)
}

func (server *GofuServer) runWithFcgi() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", server.mainHandler)
  listen, _ := net.Listen("tcp", server.bind)
  fcgi.Serve(listen, mux)
}

func (server *GofuServer) runWithHttp() {
  http.HandleFunc("/", server.mainHandler)
  http.ListenAndServe(server.bind, nil)
}

func (server *GofuServer) run() {
  if gofuConfig.Fcgi {
    server.runWithFcgi()
  } else {
    server.runWithHttp()
  }
}
