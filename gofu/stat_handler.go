package gofu

import (
  "encoding/json"
  "net/http"
  "strconv"
)

type Stat struct {
  Entries     int
  UsedStorage int
}

var StatHandler = &Handler{
  Path: "/stat",
}

func init() {
  StatHandler.Func = statHandler
}

func statHandler(w http.ResponseWriter, r *http.Request) {
  stat := &Stat{
    Entries:     cache.list.Len(),
    UsedStorage: cache.UsedStorage(),
  }

  blob, _ := json.Marshal(stat)

  w.Header().Add("Content-Type", "application/json")
  w.Header().Add("Content-Length", strconv.Itoa(len(blob)))
  w.Write(blob)
}
