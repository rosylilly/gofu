package gofu

import (
<<<<<<< HEAD
  "fmt"
  "time"
)

func bench(label string, f func()) {
  n := time.Now()
  f()
  fmt.Println(label, ":", time.Now().Sub(n))
=======
  "time"
  "fmt"
)

func bench(label string, f func()) {
  now := time.Now()
  f()
  fmt.Println(label, ":", time.Now().Sub(now))
>>>>>>> Fixes
}
