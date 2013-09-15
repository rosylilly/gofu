package main

import (
  "github.com/rosylilly/gofu/gofu"
  "fmt"
  "time"
  "runtime"
)

func main() {
  runtime.GOMAXPROCS(runtime.NumCPU())

  config := gofu.NewConfig()
  config.Load("config.json")

  go func() {
    fmt.Println(runtime.NumCPU())
    for {
      fmt.Println("NumGoroutine = ", runtime.NumGoroutine())
      time.Sleep( 1 * time.Second )
    }
  }()

  server := gofu.NewServer(config)
  server.Start()
}
