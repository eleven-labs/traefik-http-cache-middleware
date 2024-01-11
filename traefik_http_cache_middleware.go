package traefikhttpcachemiddleware

import (
  "context"
  "fmt"
  "net/http"

  "github.com/redis/go-redis/v9"
)

type Config struct {
  Redis   RedisConfig   `json:"redis" yaml:"redis" toml:"redis"`
  Paths   string        `json:"paths" yaml:"paths" toml:"paths"`
}

type RedisConfig struct {
  Hostname  string  `json:"hostname" yaml:"hostname" toml:"hostname"`
  Port      uint16  `json:"port" yaml:"port" toml:"port"`
  User      string  `json:"user" yaml:"user" toml:"user"`
  Password  string  `json:"password" yaml:"password" toml:"password"`
  Database  uint8   `json:"database" yaml:"database" toml:"database"`
  Protocol  uint8   `json:"protocol" yaml:"protocol" toml:"protocol"`
}

func CreateConfig() *Config {
	return &Config{
    Redis: RedisConfig{
      Hostname:     "localhost",
      Port:         6379,
      User:         "",
      Password:     "",
      Database:     0,
      Protocol:     3,
    },
    Paths:  "paths",
  }
}

type HttpCache struct {
  next    http.Handler
  name    string
}

func NewRedisClient(redisConfig *RedisConfig) *redis.Client {
  url := fmt.Sprintf(
    "redis://%s:%s@%s:%d/%d?protocol=%d", 
    redisConfig.User,
    redisConfig.Password,
    redisConfig.Hostname,
    redisConfig.Port,
    redisConfig.Database,
    redisConfig.Protocol,
  )

  opts, err := redis.ParseURL(url)

  if err != nil {
    panic(err)
  }

  return redis.NewClient(opts)
}


func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
  return &HttpCache {
    next,
    name,
  }, nil
}

func (hc *HttpCache) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
  rw.Header().Add("Test-Cache-header", "OK")
  hc.next.ServeHTTP(rw, req)
}

