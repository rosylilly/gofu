package gofu

import "time"

const (
  DefaultExpire = 5 * time.Minute
)

type CacheConfig struct {
  Dir            string `json:"dir"`
  Expire         string `json:"expire"`
  MaxStorageSize int    `json:"max_storage_size"`
}

func NewCacheConfig() CacheConfig {
  return CacheConfig{
    Dir:            "tmp/cache",
    Expire:         "3m",
    MaxStorageSize: 100 * 1024 * 1024,
  }
}

func (c *CacheConfig) ExpireTime() time.Duration {
  t, err := time.ParseDuration(c.Expire)
  if err != nil {
    return DefaultExpire
  }
  return t
}
