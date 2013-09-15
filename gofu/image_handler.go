package gofu

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "strconv"
)

var ImageHandler = &Handler{
  Path: "/i/",
}

func init() {
  ImageHandler.Func = imageHandler
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
  p := r.URL.Path[3:]
  fmt.Println(p)

  image, err := GetImage(p)
  if err != nil {
    fmt.Println(err)
    return
  }

  image.Blob, err = ioutil.ReadFile(image.Path)
  if err != nil {
    fmt.Println(err)
    return
  }

  w.WriteHeader(http.StatusOK)
  w.Header().Add("Content-Type", http.DetectContentType(image.Blob))
  w.Header().Add("Content-Length", strconv.Itoa(len(image.Blob)))
  w.Write(image.Blob)
}
