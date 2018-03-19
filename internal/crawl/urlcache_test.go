package crawl

import (
	"net/url"
	"testing"
)

func TestUrlCache(t *testing.T) {
	c := newUrlCache()
	u, _ := url.Parse("http://www.google.com")

	stored := c.addToCacheIfNotExists(u)
	if !stored {
		t.Error("Url could not be stored in empty cache")
	}

	_, ok := c.crawledUrls[*u]
	if !ok {
		t.Error("Url not stored in cache")
	}

	stored = c.addToCacheIfNotExists(u)
	if stored {
		t.Error("Url can't be stored twice")
	}
}
