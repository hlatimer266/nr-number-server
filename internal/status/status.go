package status

import (
	"context"
	"fmt"
	"time"

	"github.com/hlatimer266/nr-number-server/internal/cache"
)

func ReportLatest(ctx context.Context, c *cache.NumberCache) {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			summary(c)
		}
	}
}

func summary(c *cache.NumberCache) {
	fmt.Printf("Received %v unique numbers, %v duplicates. Unique total: %v\n", c.LatestUniqueNumbers.ItemCount(), c.NumberDedupes, c.AllUniqueNumbers.ItemCount())
}
