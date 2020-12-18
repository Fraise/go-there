package cache

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-there/config"
	"go-there/data"
	"strconv"
	"time"

	rediscache "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type Cache struct {
	rc *rediscache.Cache
}

// Init initializes the Redis cache from the configuration. Returns nil if the cache is not enabled.
func Init(config *config.Configuration) *Cache {
	if !config.Cache.Enabled {
		return nil
	}

	cache := new(Cache)

	client := redis.NewClient(&redis.Options{
		Addr:     config.Cache.Address + ":" + strconv.Itoa(config.Cache.Port),
		Username: config.Cache.User,
		Password: config.Cache.Password,
	})

	cache.rc = rediscache.New(&rediscache.Options{
		Redis:      client,
		LocalCache: rediscache.NewTinyLFU(1000, time.Minute),
	})

	return cache
}

// GetTarget gets a target in the cache from a path. Returns a data.ErrRedis if it fails. Returns "", nil if no cache
// exists.
func (cache *Cache) GetTarget(c *gin.Context, path string) (string, error) {
	if cache == nil {
		return "", nil
	}

	var target string
	err := cache.rc.Get(c, path, &target)

	if err != nil {
		return "", fmt.Errorf("%w: %s", data.ErrRedis, err)
	}

	return target, nil
}

// AddTarget adds a target to the cache with a ttl of 1 hour. Returns a data.ErrRedis if it fails. Returns
// nil if no cache exists.
func (cache *Cache) AddTarget(c *gin.Context, path data.Path) error {
	if cache == nil {
		return nil
	}

	err := cache.rc.Set(&rediscache.Item{
		Ctx:   c,
		Key:   path.Path,
		Value: path.Target,
		TTL:   time.Hour,
	})

	if err != nil {
		return fmt.Errorf("%w: %s", data.ErrRedis, err)
	}

	return nil
}

// DeleteTargets deletes all targets corresponding to the paths array provided. Returns a data.ErrRedis if it fails.
// Returns nil if no cache exists.
func (cache *Cache) DeleteTargets(c *gin.Context, paths []string) error {
	if cache == nil {
		return nil
	}

	var err error

	for _, p := range paths {
		if cacheErr := cache.rc.Delete(c, p); cacheErr != nil {
			err = cacheErr
		}
	}

	if err != nil {
		return fmt.Errorf("%w: %s", data.ErrRedis, err)
	}

	return nil
}
