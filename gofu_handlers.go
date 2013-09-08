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

func (server *GofuServer) getImage(reqPath string) (p string, e error) {
  cachePath := path.Join(gofuConfig.Dir.Cache, reqPath)

  fileInfo, err := os.Stat(cachePath)
  if err != nil {
    blob, err := server.bucket.Get(reqPath)

    if err != nil {
      return "", err
    }

    os.MkdirAll(path.Dir(cachePath), 0700)

    io, _ := os.Create(cachePath)
    defer io.Close()
    io.Write(blob)
  } else {
    fmt.Println(fileInfo.ModTime().Unix())
  }

  return cachePath, err
}

func (server *GofuServer) processImage(mw *imagick.MagickWand, filePath string, query url.Values) {
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

  mw.SetOption("jpeg:size", fmt.Sprintf("%dx%d", width, height))

  bench("Read Image", func() {
    mw.ReadImage(filePath)
  })

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
  var filePath string
  var err error

  bench("Load by S3 or Cache", func() {
    filePath, err = server.getImage(req.URL.Path[1:])
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
    server.processImage(magickWand, filePath, req.URL.Query())
  })

  bench("Write Response", func() {
    res.Body = magickWand.GetImageBlob()
  })
}
