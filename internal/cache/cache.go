package cache

import (
	"sync"

	goc "github.com/patrickmn/go-cache"
)

type NumberCache struct {
	AllUniqueNumbers    *goc.Cache
	LatestUniqueNumbers *goc.Cache
	NumberDedupes       int
	sync.RWMutex
}

func NewNumberCache() *NumberCache {
	return &NumberCache{
		AllUniqueNumbers:    goc.New(goc.NoExpiration, goc.NoExpiration),
		LatestUniqueNumbers: goc.New(goc.NoExpiration, goc.NoExpiration),
		NumberDedupes:       0,
	}
}

func (nc *NumberCache) AddDedupe() {
	nc.Lock()
	nc.NumberDedupes++
	nc.Unlock()
}

func (nc *NumberCache) Reset() {
	nc.Lock()
	nc.NumberDedupes = 0
	nc.Unlock()
}
