package main

import (
  "flag"
  "github.com/gographics/imagick/imagick"
  "os"
  "runtime"
)

func init() {
  gofuConfig.init()
}

func main() {
  runtime.GOMAXPROCS(runtime.NumCPU())

  flag.Parse()

  if gofuConfig.Usage {
    flag.Usage()
    os.Exit(0)
  }

  imagick.Initialize()
  defer func() {
    imagick.Terminate()
  }()

  e := gofuConfig.load()
  if e != nil {
    os.Exit(1)
  }

  start()
}
