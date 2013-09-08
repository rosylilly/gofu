package main

import (
  "fmt"
  "github.com/gographics/imagick/imagick"
  "net/http"
  "net/url"
  "strconv"
  "strings"
  "path"
  "time"
  "os"
)

type CropContext struct {
  width uint
  height uint
  x int
  y int
}

type RequestContext struct {
  mtime int64
  width uint
  height uint
  quarity uint
  blur float64
  crop *CropContext
}

func (ctx *RequestContext) isResize() bool {
  return (ctx.width > 0 && ctx.height > 0)
}

func bench(label string, f func()) {
  t := time.Now().UnixNano()
  for(len(label) < 30) {
    label += " "
  }
  f()
  if gofuConfig.Verbose {
    fmt.Printf("%s: %d\n", label, (time.Now().UnixNano() - t)/1000000)
  }
}

func atoi64(query string, source int64) int64 {
  to, err := strconv.ParseInt(query, 10, 32)
  if err == nil {
    source = to
  }
  return source
}

func atoi(query string, source int) int {
  to, err := strconv.ParseInt(query, 10, 16)
  if err == nil {
    source = int(to)
  }
  return source
}

func atoui(query string, source uint) uint {
  to, err := strconv.ParseUint(query, 10, 16)
  if err == nil {
    source = uint(to)
  }
  return source
}

func atof(query string, source float64) float64 {
  to, err := strconv.ParseFloat(query, 64)
  if err == nil {
    source = float64(to)
  }
  return source
}

func (server *GofuServer) getImage(reqPath string, context *RequestContext) (p string, e error) {
  cachePath := path.Join(gofuConfig.Dir.Cache, reqPath)

  fileInfo, err := os.Stat(cachePath)
  if err != nil || fileInfo.ModTime().Unix() < context.mtime {
    blob, err := server.bucket.Get(reqPath)

    if err != nil {
      return "", err
    }

    os.MkdirAll(path.Dir(cachePath), 0700)

    io, _ := os.Create(cachePath)
    defer io.Close()
    io.Write(blob)
  }

  return cachePath, err
}

func (server *GofuServer) parseQuery(query url.Values, context *RequestContext) {
  context.blur = gofuConfig.Image.DefaultBlur
  context.quarity = gofuConfig.Image.DefaultQuarity

  if query["m"] != nil {
    context.mtime = atoi64(query["m"][0], context.mtime)
  }
  if query["w"] != nil {
    context.width = atoui(query["w"][0], context.width)
  }
  if query["h"] != nil {
    context.height = atoui(query["h"][0], context.height)
  }
  if query["q"] != nil {
    context.quarity = atoui(query["q"][0], context.quarity)
  }
  if query["b"] != nil {
    context.blur = atof(query["b"][0], context.blur)
  }
  if query["c"] != nil {
    cropQuery := strings.Split(query["c"][0], ",")
    if len(cropQuery[0]) > 0 && len(cropQuery[1]) > 0 && len(cropQuery[2]) > 0 && len(cropQuery[3]) > 0 {
      context.crop = new(CropContext)
      context.crop.width = atoui(cropQuery[0], context.width)
      context.crop.height = atoui(cropQuery[1], context.height)
      context.crop.x = atoi(cropQuery[2], 0)
      context.crop.y = atoi(cropQuery[3], 0)
    }
  }
}

func (server *GofuServer) processImage(mw *imagick.MagickWand, filePath string, ctx *RequestContext) {
  if ctx.isResize() {
    mw.SetOption("jpeg:size", fmt.Sprintf("%dx%d", ctx.width, ctx.height))
  }

  bench("Read Image", func() {
    mw.ReadImage(filePath)
  })

  if ctx.crop != nil {
    crop := ctx.crop
    mw.CropImage(crop.width, crop.height, crop.x, crop.y)
  }

  if ctx.isResize() {
    mw.ResizeImage(ctx.width, ctx.height, imagick.FILTER_CUBIC, ctx.blur)
  }

  mw.SetImageCompressionQuality(ctx.quarity)
  mw.StripImage()
}

func (server *GofuServer) imageHandler(res *GofuResponse, req *http.Request) {
  var filePath string
  var err error

  ctx := new(RequestContext)
  server.parseQuery(req.URL.Query(), ctx)

  bench("Load by S3 or Cache", func() {
    filePath, err = server.getImage(req.URL.Path[1:], ctx)
    if err != nil {
      res.Status = http.StatusNotFound
      res.ClearBody()
      return
    }
  })

  var magickWand *imagick.MagickWand
  bench("Create MagickWand", func() {
    magickWand = <-server.wands
  })
  defer func() {
    magickWand.Clear()
    server.wands <- magickWand
  }()

  bench("ProcessImage", func() {
    server.processImage(magickWand, filePath, ctx)
  })

  bench("Write Response", func() {
    res.Body = magickWand.GetImageBlob()
  })
}
