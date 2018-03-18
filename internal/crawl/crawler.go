package crawl

import (
	"fmt"
	"net/url"
)

// Choose this to be basically infinite.
// We want to avoid blocking the worker goroutines.
// Alternatively implement "stacked" channel or
// custom queue with synchronisation points
const urlFrontierBufferSize int = 1e6

type Crawler struct {
	// Filtered URLs to be crawled
	urlFrontier chan *url.URL
	// Responses from individual workers
	fetchedURLs chan []*url.URL
	// Cache for visited URLs
	urlCache urlCache
	// Keeping track of worker tasks
	tasksInProgress int
}

func (c *Crawler) Run(seedUrl url.URL, numWorkers int) {
	fmt.Println("Crawling", seedUrl.String())
	c.urlFrontier = make(chan *url.URL, urlFrontierBufferSize)
	c.fetchedURLs = make(chan []*url.URL, numWorkers)
	c.urlCache = *newUrlCache()
	c.launchWorkerPool(numWorkers)
	c.enqueueUrls([]*url.URL{&seedUrl})
	for c.tasksInProgress > 0 {
		select {
		case fetchedUrls := <-c.fetchedURLs:
			c.tasksInProgress--
			c.enqueueUrls(fetchedUrls)
		}
	}
	fmt.Println("Done crawling")
}

func (c *Crawler) launchWorkerPool(numWorkers int) {
	// We could launch the workers dynamically,
	// but since we expect full load on them after
	// an initial warm-up period, we can set up the pool directly
	fmt.Println("Starting up workers...")
	for i := 0; i < numWorkers; i++ {
		w := worker{
			fetchedURLs: c.fetchedURLs,
			urlFrontier: c.urlFrontier,
		}
		go w.Run()
	}
}

func (c *Crawler) enqueueUrls(urlsToEnqueue []*url.URL) {
	for _, urlToEnqueue := range urlsToEnqueue {
		if isNewUrl := c.urlCache.addToCacheIfNotExists(urlToEnqueue); isNewUrl {
			select {
			case c.urlFrontier <- urlToEnqueue:
				// url is queued for processing
				c.tasksInProgress++
			default:
				// Our queue is overloaded. Let's drop the url.
				// Alternatively implement infinite queue without backpressure
				fmt.Println("URL frontier overloaded. Dropping url", urlToEnqueue)
			}
		}
	}
}
