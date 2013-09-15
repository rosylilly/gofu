package gofu

import "time"

const (
  DefaultTimeout = 5 * time.Second
)

type TimeoutConfig struct {
  Read  string `json:"read"`
  Write string `json:"write"`
}

func NewTimeoutConfig() TimeoutConfig {
  return TimeoutConfig{
    Read:  "5m",
    Write: "5m",
  }
}

func (c *TimeoutConfig) ReadTime() time.Duration {
  t, err := time.ParseDuration(c.Read)
  if err != nil {
    return DefaultTimeout
  }
  return t
}

func (c *TimeoutConfig) WriteTime() time.Duration {
  t, err := time.ParseDuration(c.Write)
  if err != nil {
    return DefaultTimeout
  }
  return t
}
