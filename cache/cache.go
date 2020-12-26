package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
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

	// Configure local cache
	var localCache rediscache.LocalCache

	if config.Cache.LocalCacheSize <= 0 || config.Cache.LocalCacheTtlSec <= 0 {
		localCache = nil
		log.Warn().Msg("cache enabled, but no local cache configured")
	} else {
		localCache = rediscache.NewTinyLFU(
			config.Cache.LocalCacheSize,
			time.Second*time.Duration(config.Cache.LocalCacheTtlSec),
		)
	}

	// Configure network cache
	cache := new(Cache)

	// Never retries if it cannot connect to the instance. It will still tries to connect for each request, but it
	// prevents the total request time to be super long (because of multiple retries) if if fails.
	client := redis.NewClient(&redis.Options{
		Network:    "",
		Addr:       config.Cache.Address + ":" + strconv.Itoa(config.Cache.Port),
		Username:   config.Cache.User,
		Password:   config.Cache.Password,
		MaxRetries: -1,
	})

	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		log.Error().Err(fmt.Errorf("%w: %s", data.ErrRedis, err)).
			Msg("cannot ping the configured redis instance, using local cache only")
	}

	// use a local cache of 1000 elements
	cache.rc = rediscache.New(&rediscache.Options{
		Redis:      client,
		LocalCache: localCache,
	})

	return cache
}

// GetTarget gets a target in the cache from a path. Returns a data.ErrRedis if it fails. Returns "", nil if no cache
// exists.
func (cache *Cache) GetTarget(path string) (string, error) {
	if cache == nil {
		return "", nil
	}

	var target string
	err := cache.rc.Get(context.Background(), path, &target)

	if err != nil {
		if !errors.Is(err, rediscache.ErrCacheMiss) {
			return "", fmt.Errorf("%w: %s", data.ErrRedis, err)
		}
	}

	return target, nil
}

// AddTarget adds a target to the cache with a ttl of 1 hour. Returns a data.ErrRedis if it fails. Returns
// nil if no cache exists.
func (cache *Cache) AddTarget(path data.Path) error {
	if cache == nil {
		return nil
	}

	err := cache.rc.Set(&rediscache.Item{
		Ctx:   context.Background(),
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
func (cache *Cache) DeleteTargets(paths []string) error {
	if cache == nil {
		return nil
	}

	var err error

	for _, p := range paths {
		if cacheErr := cache.rc.Delete(context.Background(), p); cacheErr != nil {
			err = cacheErr
		}
	}

	if err != nil {
		return fmt.Errorf("%w: %s", data.ErrRedis, err)
	}

	return nil
}
