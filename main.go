package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "github.com/gographics/imagick/imagick"
  "io/ioutil"
  "launchpad.net/goamz/aws"
  "os"
  "runtime"
)

type Config struct {
  Path     string
  Help     bool
  Bind     string
  Port     uint
  Bucket   string
  Fcgi     bool
  MaxCache int
  S3       aws.Auth
}

var gofu_config Config

func init() {
  gofu_config.Bind = ""
  gofu_config.Port = 8808
  gofu_config.MaxCache = 100
  gofu_config.Fcgi = false

  flag.StringVar(&gofu_config.Path, "c", "./config.json", "config file path")
  flag.BoolVar(&gofu_config.Help, "h", false, "show help")

  flag.Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage of gofu:\n")
    flag.PrintDefaults()
  }
}

func main() {
  runtime.GOMAXPROCS(runtime.NumCPU())

  flag.Parse()

  if gofu_config.Help {
    flag.Usage()
    os.Exit(0)
  }

  imagick.Initialize()
  defer func() {
    imagick.Terminate()
  }()

  file, e := ioutil.ReadFile(gofu_config.Path)
  if e != nil {
    fmt.Printf("Config file error: %v\n", e)
    os.Exit(1)
  }

  json.Unmarshal(file, &gofu_config)

  start()
}
