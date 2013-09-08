package main

import (
  "net/http"
  "strconv"
)

type GofuResponse struct {
  writer http.ResponseWriter
  Status int
  Head   bool

  Body []byte
}

func NewGofuResponse(writer http.ResponseWriter) *GofuResponse {
  response := new(GofuResponse)

  response.writer = writer
  response.Status = http.StatusOK

  return response
}

func (response *GofuResponse) AddHeader(key, val string) {
  response.writer.Header().Add(key, val)
}

func (response *GofuResponse) ClearBody() {
  response.Body = []byte("")
}

func (response *GofuResponse) Write() {
  w := response.writer

  w.WriteHeader(response.Status)

  w.Header().Add("Content-Type", http.DetectContentType(response.Body))
  w.Header().Add("Content-Length", strconv.Itoa(len(response.Body)))

  if !response.Head {
    w.Write(response.Body)
  }
}
