package crawl

import "net/url"

type urlCache struct {
	crawledUrls map[url.URL]struct{}
}

func newUrlCache() *urlCache {
	c := new(urlCache)
	c.crawledUrls = make(map[url.URL]struct{})
	return c
}

func (c *urlCache) addToCacheIfNotExists(urlToCache url.URL) bool {
	// Adds url to cache and returns a bool whether this was added
	// (true) or whether it already existed in the cache (false)
	_, ok := c.crawledUrls[urlToCache]
	if !ok {
		c.crawledUrls[urlToCache] = struct{}{}
		return true
	}
	return false
}
