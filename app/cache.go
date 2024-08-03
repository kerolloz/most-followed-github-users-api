package main

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var cacheContainer = cache.New(5*time.Minute, 10*time.Minute)

// The following function tries to get the key from the cache. If the key is not found, it calls the provided function to get the value and stores it in the cache.
func GetFromCacheOrEvaluateFunction(key string, f func() interface{}) (interface{}, bool) {
	if value, found := cacheContainer.Get(key); found {
		return value, true
	}

	value := f()
	cacheContainer.SetDefault(key, value)
	return value, false
}
