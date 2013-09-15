package gofu

import (
  "container/list"
  "crypto/md5"
  "fmt"
  "io"
  "os"
  "path"
  "strings"
  "sync"
  "time"
)

const (
  separateStep = 4
)

type Cache struct {
  mutex   sync.Mutex
  dir     string
  expire  time.Duration
  maxSize int
  list    *list.List
  entries map[string]*list.Element
}

type entry struct {
  key  string
  path string
  size int
}

var cache *Cache

func InitCache(c CacheConfig) {
  cache = NewCache(c)
}

func NewCache(config CacheConfig) *Cache {
  return &Cache{
    dir:     config.Dir,
    expire:  config.ExpireTime(),
    maxSize: config.MaxStorageSize,
    list:    list.New(),
    entries: make(map[string]*list.Element),
  }
}

func (c *Cache) PathByKey(key string) string {
  hash := md5.New()
  io.WriteString(hash, key)
  hashKey := fmt.Sprintf("%x", hash.Sum(nil))
  generateKey := []string{}
  for i, l := 0, (len(hashKey) / separateStep); i < l; i++ {
    index := i * separateStep
    generateKey = append(generateKey, string(hashKey[index:index+separateStep]))
  }
  return path.Join(c.dir, strings.Join(generateKey, "/"))
}

func (c *Cache) Fetch(
  key string,
  missing func(string) ([]byte, error),
) (string, error) {
  path, ok := c.Get(key)

  if !ok {
    blob, err := missing(key)
    if err != nil {
      return "", err
    }
    c.Set(key, blob)
    path = c.PathByKey(key)
  }

  return path, nil
}

func (c *Cache) Get(key string) (string, bool) {
  if element, ok := c.entries[key]; ok {
    c.list.MoveToFront(element)
    return element.Value.(*entry).path, true
  }
  return "", false
}

func (c *Cache) Set(key string, blob []byte) error {
  var cachePath string
  var ok bool
  if cachePath, ok = c.Get(key); !ok {
    cachePath = c.PathByKey(key)
  }

  err := os.MkdirAll(path.Dir(cachePath), 0700)
  if err != nil {
    return err
  }

  file, err := os.Create(cachePath)
  defer file.Close()

  if err != nil {
    return err
  }

  bytes, err := file.Write(blob)
  if err != nil {
    return err
  }

  var element *list.Element
  var hit bool
  if element, hit = c.entries[key]; !hit {
    element = c.list.PushFront(&entry{key, cachePath, bytes})
  } else {
    element.Value.(*entry).size = bytes
  }
  c.entries[key] = element

  for c.UsedStorage() > c.maxSize {
    c.RemoveOldest()
  }

  return nil
}

func (c *Cache) UsedStorage() int {
  var bytes int
  for element := c.list.Front(); element != nil; element = element.Next() {
    bytes = bytes + element.Value.(*entry).size
  }
  return bytes
}

func (c *Cache) Remove(key string) {
  if element, hit := c.entries[key]; hit {
    c.removeElement(element)
  }
}

func (c *Cache) RemoveOldest() {
  element := c.list.Back()
  if element != nil {
    c.removeElement(element)
  }
}

func (c *Cache) removeElement(element *list.Element) {
  c.list.Remove(element)
  entry := element.Value.(*entry)
  os.Remove(entry.path)
  delete(c.entries, entry.key)
}
