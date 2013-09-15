package gofu

import (
  "encoding/json"
  "io/ioutil"
)

type Config struct {
  Listen   string        `json:"listen"`
  FcgiMode bool          `json:"fcgi_mode"`
  PidFile  string        `json:"pid_file"`
  Timeout  TimeoutConfig `json:"timeout"`
  Cache    CacheConfig   `json:"cache"`
  S3       S3Config      `json:"s3"`

  path string
}

func NewConfig() *Config {
  return &Config{}
}

func (c *Config) Load(path string) error {
  jsonFile, err := ioutil.ReadFile(path)
  if err != nil {
    return err
  }

  return json.Unmarshal(jsonFile, &c)
}
