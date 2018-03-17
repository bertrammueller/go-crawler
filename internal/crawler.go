package crawl

import (
	"fmt"
	"net/url"
)

const maxNumWorkers int = 1
const urlFrontierBufferSize int = 1e6

type Crawler struct {
	// Filtered URLs to be crawled
	urlFrontier chan url.URL
	// Responses from individual workers
	fetchedURLs chan url.URL
	// Cache for visited URLs
	urlCache urlCache
}

func (c *Crawler) Run(seedUrl url.URL) {
	fmt.Println("Crawling", seedUrl)
	c.urlFrontier = make(chan url.URL, urlFrontierBufferSize)
	c.fetchedURLs = make(chan url.URL)
	c.urlCache = *newUrlCache()
	c.launchWorkerPool(maxNumWorkers)
	c.enqueueUrl(seedUrl)
	for {
		select {
		case fetchedUrl := <-c.fetchedURLs:
			c.enqueueUrl(fetchedUrl)
		}
	}
}

func (c *Crawler) launchWorkerPool(numWorkers int) {
	// We could launch the workers dynamically,
	// but since we expect full load on them after
	// an initial warm-up period, we can set up the pool directly
	for i := 0; i < numWorkers; i++ {
		w := worker{
			fetchedURLs: c.fetchedURLs,
			urlFrontier: c.urlFrontier,
		}
		go w.Run()
	}
}

func (c *Crawler) enqueueUrl(urlToEnqueue url.URL) {
	if isNewUrl := c.urlCache.addToCacheIfNotExists(urlToEnqueue); isNewUrl {
		select {
		case c.urlFrontier <- urlToEnqueue:
			// url is queued for processing
		default:
			// Our queue is overloaded. Let's drop the url.
			// Alternatively implement infinite queue without backpressure
			fmt.Println("URL frontier overloaded. Dropping url", urlToEnqueue)
		}
	}
}
