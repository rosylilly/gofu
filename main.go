package main

import (
  "github.com/rosylilly/gofu/gofu"
  "runtime"
)

func main() {
  runtime.GOMAXPROCS(runtime.NumCPU())

  config := gofu.NewConfig()
  config.Load("config.json")

  gofu.InitS3Client(config.S3)
  gofu.InitCache(config.Cache)

  server := gofu.NewServer(config)
  server.Start()
}
