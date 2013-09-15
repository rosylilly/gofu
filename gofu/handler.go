package gofu

import (
  "net/http"
)

type Handler struct {
  Path string
  Func func(http.ResponseWriter, *http.Request)
}
