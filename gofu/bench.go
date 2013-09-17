package gofu

import (
  "fmt"
  "time"
)

func bench(label string, f func()) {
  n := time.Now()
  f()
  fmt.Println(label, ":", time.Now().Sub(n))
}
