package gofu

import (
  "time"
  "fmt"
)

func bench(label string, f func()) {
  now := time.Now()
  f()
  fmt.Println(label, ":", time.Now().Sub(now))
}
