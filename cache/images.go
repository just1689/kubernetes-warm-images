package cache

import (
	"github.com/ReneKroon/ttlcache/v2"
	"time"
)

var (
	c        ttlcache.SimpleCache = initCache()
	notFound                      = ttlcache.ErrNotFound
)

func initCache() *ttlcache.Cache {
	result := ttlcache.NewCache()
	//TODO: move to ENV... and then to chart
	result.SetTTL(5 * time.Minute)

	return result
}

func ShouldPullImage(image string) (proceed bool) {
	_, err := c.Get(image)
	if err == notFound {
		proceed = true
		c.Set(image, true)
		return
	}
	return false

}
