package gofu

import (
  "net/http"
)

type Server struct {
  http *http.Server
  mux  *http.ServeMux
}

func NewServer(config *Config) *Server {
  serveMux := http.NewServeMux()

  httpServer := &http.Server{
    Addr:         config.Listen,
    ReadTimeout:  config.Timeout.ReadTime(),
    WriteTimeout: config.Timeout.WriteTime(),
    Handler:      serveMux,
  }

  server := &Server{
    http: httpServer,
    mux:  serveMux,
  }

  server.init()

  return server
}

func (s *Server) Start() error {
  return s.http.ListenAndServe()
}

func (s *Server) AddHandler(handler *Handler) {
  s.mux.HandleFunc(handler.Path, handler.Func)
}

func (s *Server) init() {
  s.AddHandler(HealthHandler)
}
