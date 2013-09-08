package main

import (
  "github.com/gographics/imagick/imagick"
  "net/http"
  "net/url"
  "strconv"
  "strings"
)

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

func (server *GofuServer) getImage(path string) (b []byte, e error) {
  byCache, r := server.lru.Get(path)

  if !r {
    blob, err := server.bucket.Get(path)

    if err == nil {
      server.lru.Add(path, blob)
    }

    return blob, err
  }

  return byCache.([]byte), nil
}

func (server *GofuServer) processImage(mw *imagick.MagickWand, query url.Values) {
  originWidth := mw.GetImageWidth()
  originHeight := mw.GetImageHeight()
  width := originWidth
  height := originHeight
  blur := gofuConfig.Image.DefaultBlur
  quarity := gofuConfig.Image.DefaultQuarity

  if query["w"] != nil {
    width = atoui(query["w"][0], width)
  }
  if query["h"] != nil {
    height = atoui(query["h"][0], height)
  }
  if query["q"] != nil {
    quarity = atoui(query["q"][0], quarity)
  }
  if query["b"] != nil {
    blur = atof(query["b"][0], blur)
  }
  if query["c"] != nil {
    cropQuery := strings.Split(query["c"][0], ",")
    if len(cropQuery[0]) > 0 && len(cropQuery[1]) > 0 && len(cropQuery[2]) > 0 && len(cropQuery[3]) > 0 {
      cropWidth := atoui(cropQuery[0], width)
      cropHeight := atoui(cropQuery[1], height)
      cropPointX := atoi(cropQuery[2], 0)
      cropPointY := atoi(cropQuery[3], 0)
      mw.CropImage(cropWidth, cropHeight, cropPointX, cropPointY)
    }
  }

  if width != originWidth || height != originHeight {
    mw.ResizeImage(width, height, imagick.FILTER_CUBIC, blur)
  }

  mw.SetImageCompressionQuality(quarity)
  mw.StripImage()
}

func (server *GofuServer) imageHandler(res *GofuResponse, req *http.Request) {
  blob, err := server.getImage(req.URL.Path[1:])
  if err != nil {
    res.Status = http.StatusNotFound
    res.ClearBody()
    return
  }

  magickWand := <-server.wands
  defer func() {
    magickWand.Clear()
    server.wands <- magickWand
  }()
  magickWand.ReadImageBlob(blob)

  server.processImage(magickWand, req.URL.Query())

  res.Body = magickWand.GetImageBlob()
}
