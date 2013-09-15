package gofu

import (
  "net/http"
  "os"
  "os/signal"
  "syscall"
)

const (
  SERVICE_OUT = "NG"
  SERVICE_IN  = "OK"
)

var HealthHandler = &Handler{
  Path: "/health",
}

var alive = true

func init() {
  HealthHandler.Func = healthHandler
  go watchSignal()
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Add("Content-Type", "text/plain")
  w.Header().Add("Content-Length", "2")

  if alive {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(SERVICE_IN))
  } else {
    w.WriteHeader(http.StatusServiceUnavailable)
    w.Write([]byte(SERVICE_OUT))
  }
}

func watchSignal() {
  for {
    channel := make(chan os.Signal, 1)
    signal.Notify(channel, syscall.SIGWINCH)

    <-channel
    alive = !alive
  }
}
