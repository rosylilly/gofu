package gofu

import "time"

const (
  DefaultExpire = 5 * time.Minute
)

type CacheConfig struct {
  Dir            string `json:"dir"`
  Expire         string `json:"expire"`
  MaxStorageSize int64  `json:"max_storage_size"`
}

func (c *CacheConfig) ExpireTime() time.Duration {
  t, err := time.ParseDuration(c.Expire)
  if err != nil {
    return DefaultExpire
  }
  return t
}
