package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "io/ioutil"
  "launchpad.net/goamz/aws"
)

type GofuConfig struct {
  Path string
  Usage bool
  Bind string
  Port uint
  Bucket string
  Fcgi bool
  MaxCache int
  S3Config aws.Auth
  Verbose bool
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
  config.MaxCache = 1000
  config.Fcgi = false
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
