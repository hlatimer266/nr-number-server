package write

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hlatimer266/nr-number-server/internal/cache"
	goc "github.com/patrickmn/go-cache"
)

const (
	path          = "number.log"
	reqLength     = 9
	serverNewLine = "$"
)

func Latest(ctx context.Context, c *cache.NumberCache) (err error) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err = numFile(c)
			if err != nil {
				return
			}
		}
	}
}

func numFile(c *cache.NumberCache) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for k := range c.LatestUniqueNumbers.Items() {
		c.AllUniqueNumbers.Set(k, 1, goc.NoExpiration)
		if _, err = f.WriteString(k + "\n"); err != nil {
			return err
		}
		c.LatestUniqueNumbers.Delete(k)
	}

	c.Reset()
	return nil
}

func isValidInput(num string) bool {
	// checks input has server defined new lines, is an integer and is 9 characters long
	nl := num[len(num)-1:]
	_, err := strconv.Atoi(num[0 : len(num)-1])
	if nl != serverNewLine || len(num[0:len(num)-1]) != reqLength || err != nil {
		return false
	}
	return true
}

func NumCache(num string, c *cache.NumberCache) error {

	if !isValidInput(num) {
		return fmt.Errorf("error: %v is invalid", num)
	}

	// check if latest number received is in cached set
	_, found := c.AllUniqueNumbers.Get(num)
	_, foundUnique := c.LatestUniqueNumbers.Get(num)
	if !found && !foundUnique {
		c.LatestUniqueNumbers.Set(num, 1, goc.NoExpiration)
	} else {
		c.AddDedupe()
	}

	return nil
}
