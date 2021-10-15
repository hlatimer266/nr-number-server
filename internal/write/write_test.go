package write

import (
	"os"
	"testing"

	"github.com/hlatimer266/nr-number-server/internal/cache"
	goc "github.com/patrickmn/go-cache"
)

func TestWrite(t *testing.T) {

	t.Run("cache can write to file and decreases LatestUniqueNumbers count", func(t *testing.T) {
		c := cache.NewNumberCache()
		c.LatestUniqueNumbers.Set("123456789$", 1, goc.NoExpiration)
		err := numFile(c)
		if err != nil {
			t.Fatal("failed to write cached entry to file")
		}
		if c.LatestUniqueNumbers.ItemCount() != 0 {
			t.Fatal("failed to remove number from cache")
		}
		os.Remove("number.log")
	})

	t.Run("test valid and invalid input", func(t *testing.T) {
		var cases = map[string]bool{
			"123456789$": true,
			"123456789":  false,
			"12345":      false,
			"blah":       false,
		}
		for k, v := range cases {
			expected := v
			actual := isValidInput(k)
			if actual != expected {
				t.Fatalf("test failed: case [%s] expected [%v] but got [%v]\n", k, expected, actual)
			}

		}

	})

	t.Run("only unique values should be added to short term cache", func(t *testing.T) {
		testCache := cache.NewNumberCache()
		var cases = []string{
			"123456789$",
			"123456799$",
			"123456899$",
			"123456789$",
		}
		for _, c := range cases {
			NumCache(c, testCache)
		}

		if testCache.LatestUniqueNumbers.ItemCount() != 3 {
			t.Fatalf("test failed: expected [3] valid cases but got [%v]\n", testCache.LatestUniqueNumbers.ItemCount())
		}

	})

	t.Run("values that exist in long term cache shouldn't be added to short term cache", func(t *testing.T) {
		testCache := cache.NewNumberCache()
		testCache.AllUniqueNumbers.Set("123456789$", 1, goc.NoExpiration)
		var cases = []string{
			"123456789$",
			"123456799$",
			"123456899$",
			"123456789$",
		}
		for _, c := range cases {
			NumCache(c, testCache)
		}

		if testCache.LatestUniqueNumbers.ItemCount() != 2 {
			t.Fatalf("test failed: expected [2] valid cases but got [%v]\n", testCache.LatestUniqueNumbers.ItemCount())
		}
	})

}
