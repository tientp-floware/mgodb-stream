package cache

import (
	"time"

	cache "github.com/go-redis/cache/v7"
	"github.com/go-redis/redis/v7"
	"github.com/vmihailenco/msgpack/v4"

	logger "g.ghn.vn/go-common/zap-logger"
	cfg "github.com/tientp-floware/mgodb-stream/config"
)

type (
	// Redis driver cache with
	Redis struct {
		Codec  *cache.Codec
		Expire time.Duration
		Err    error
	}
)

var (
	ring = redis.NewRing(&redis.RingOptions{
		Addrs: cfg.Config.Redis,
	})
	log = logger.GetLogger("[Redis cache]")
)

// NewRedisCache new instance
func NewRedisCache() *Redis {
	rc := new(Redis)
	rc.Expire = time.Hour
	rc.Codec = &cache.Codec{
		Redis: ring,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
	return rc
}

// Set add cache
func (rc *Redis) Set(key string, object interface{}) *Redis {
	rc.Err = rc.Codec.Once(&cache.Item{
		Key:        key,
		Object:     object,
		Expiration: rc.Expire,
	})
	return rc
}

// Get cache
func (rc *Redis) Get(key string, data interface{}) *Redis {
	rc.Err = rc.Codec.Get(key, &data)
	return rc
}

// Del delete cache
func (rc *Redis) Del(key string) *Redis {
	rc.Err = rc.Codec.Delete(key)
	return rc
}

// Check show error
func (rc *Redis) Check() error {
	return rc.Err
}

// Exists reports whether object for the given key exists.
func (rc *Redis) Exists(key string) bool {
	return rc.Codec.Get(key, nil) == nil
}
