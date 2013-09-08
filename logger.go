package main

import (
  "io"
  "net/http"
  "time"
  "text/template"
)

type GofuLogger struct {
  io io.Writer
  template *template.Template
}

var logFormat = map[string]string{
  "combined": "{{.Host}} {{.Status}} {{.Method}} {{.Path}} {{.Size}} [{{.RequestTime}}] {{.ResponseTime}} \"{{.UserAgent}}\"\n",
  "ltsv": "host:{{.Host}}\tstatus:{{.Status}}\tmethod:{{.Method}}\tpath:{{.Path}}\tsize:{{.Size}}\trequestTime:{{.RequestTime}}\tresponseTime:{{.ResponseTime}}\tua:{{.UserAgent}}\n",
}

type LoggingField struct {
  Host string
  Status int
  Size int
  Method string
  Path string
  ResponseTime int64
  RequestTime time.Time
  UserAgent string
}

func NewGofuLogger(io io.Writer) *GofuLogger {
  var err error
  logger := GofuLogger{
    io: io,
  }
  logger.template, err = template.New("logger").Parse(logger.format())
  if err != nil { panic(err) }

  return &logger
}

func (logger *GofuLogger) format() string {
  if len(logFormat[gofuConfig.LogFormat]) > 0 {
    return logFormat[gofuConfig.LogFormat]
  }
  return gofuConfig.LogFormat
}

func (logger *GofuLogger) log(res *GofuResponse, req *http.Request, t time.Time) {
  responseTime := (time.Now().UnixNano() - t.UnixNano()) / 1000000
  userAgent := ""
  if req.Header["User-Agent"] != nil {
    userAgent = req.Header["User-Agent"][0]
  }

  field := LoggingField{
    Host: req.Host,
    Status: res.Status,
    Size: len(res.Body),
    Method: req.Method,
    Path: req.URL.Path + "?" + req.URL.RawQuery,
    ResponseTime: responseTime,
    RequestTime: t,
    UserAgent: userAgent,
  }
  logger.template.Execute(logger.io, field)
}
