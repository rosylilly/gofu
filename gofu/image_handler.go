package gofu

import (
  "github.com/gographics/imagick/imagick"
  "net/http"
  "runtime"
)

type RequestContexts chan *RequestContext

var requestContexts RequestContexts

var ImageHandler = &Handler{
  Path: "/i/",
}

func init() {
  imagick.Initialize()

  ImageHandler.Func = imageHandler

  num := runtime.NumCPU() * 4
  requestContexts = make(RequestContexts, num)
  for i := 0; i < num; i++ {
    requestContexts <- &RequestContext{
      Wand: imagick.NewMagickWand(),
    }
  }
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
  ctx := <-requestContexts

  ctx.Execute(w, r)

  requestContexts <- ctx
}
