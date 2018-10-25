package localcache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type CacheOptions struct {
	ExpireTime  int
	CleanupTime int
}

type LocalCacheService struct {
	cache   *gocache.Cache
	expire  time.Duration
	cleanup time.Duration
}

func (service *LocalCacheService) Initialize(cacheCfg interface{}) {
	options, ok := cacheCfg.(CacheOptions)
	if !ok {
		panic("local cache service config error!")
	}
	service.expire = time.Duration(options.ExpireTime) * time.Second
	service.cleanup = time.Duration(options.CleanupTime) * time.Second
	service.cache = gocache.New(service.expire, service.cleanup)
}
