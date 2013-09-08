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
  flag.Parse()

  if gofuConfig.Usage {
    flag.Usage()
    os.Exit(0)
  }

  runtime.GOMAXPROCS(gofuConfig.MaxProc)

  imagick.Initialize()
  defer imagick.Terminate()

  e := gofuConfig.load()
  if e != nil {
    os.Exit(1)
  }

  start()
}
