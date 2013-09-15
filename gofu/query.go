package gofu

import (
  "net/url"
  "strconv"
  "strings"
)

type Query struct {
  ResizedWidth  uint
  ResizedHeight uint
}

func stringToUint(str string, def uint) uint {
  to, err := strconv.ParseUint(str, 10, 64)
  if err == nil {
    def = uint(to)
  }
  return def
}

func ParseQuery(values url.Values) *Query {
  query := &Query{}

  query.parseSize(values.Get("s"))

  return query
}

func (q *Query) parseSize(s string) {
  query := strings.Split(s, "x")

  if len(query[0]) > 0 || len(query[1]) > 0 {
    q.ResizedWidth = stringToUint(query[0], 0)
    q.ResizedHeight = stringToUint(query[1], 0)
  }
}
