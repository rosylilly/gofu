package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "io/ioutil"
  "launchpad.net/goamz/aws"
  "runtime"
)

type ImageConfig struct {
  DefaultBlur    float64
  DefaultQuarity uint
}

type GofuConfig struct {
  Path      string
  Usage     bool
  Bind      string
  Port      uint
  Bucket    string
  Fcgi      bool
  MaxCache  int
  MaxProc   int
  LogPath   string
  LogFormat string
  Image     ImageConfig
  S3Config  aws.Auth
  Verbose   bool
}

var gofuConfig GofuConfig

func (config *GofuConfig) init() {
  config.setDefault()

  flag.StringVar(&config.Path, "c", "./config.json", "config file path")
  flag.BoolVar(&config.Usage, "h", false, "show help")
  flag.BoolVar(&config.Verbose, "v", false, "verbose mode")
}

func (config *GofuConfig) setDefault() {
  config.Bind = ""
  config.Port = 8088
  config.Fcgi = false
  config.MaxCache = 1000
  config.MaxProc = runtime.NumCPU()
  config.LogPath = ""
  config.LogFormat = "combined"
  config.Image = ImageConfig{
    DefaultBlur:    float64(1),
    DefaultQuarity: uint(95),
  }
}

func (config *GofuConfig) load() (err error) {
  file, e := ioutil.ReadFile(config.Path)
  if e != nil {
    fmt.Printf("Config file error: %v\n", e)
    return e
  }

  json.Unmarshal(file, &config)

  if config.Verbose {
    fmt.Printf("%+v\n", config)
  }

  return
}
