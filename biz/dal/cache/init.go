package cache

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

var Cache *bigcache.BigCache
var err error

func Init(ctx context.Context) {
	Cache, err = bigcache.New(ctx, bigcache.DefaultConfig(60*time.Minute))
	if err != nil {
		panic(err)
	}
}
