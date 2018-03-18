package crawl

import (
	"fmt"
	"net/url"
	"time"
)

// Choose this to be basically infinite.
// We want to avoid blocking the worker goroutines.
// Alternatively implement "stacked" channel or
// custom queue with synchronisation points
const urlFrontierBufferSize int = 1e6

type workerResponse struct {
	base       *url.URL
	linkedUrls []*url.URL
}

type Crawler struct {
	// Filtered URLs to be crawled
	urlFrontier chan *url.URL
	// Responses from individual workers
	workerResp chan workerResponse
	// Cache for visited URLs
	urlCache urlCache
	// Keeping track of worker tasks
	tasksInProgress int
	// Continuously updated sitemap
	sitemap sitemap
}

func (c *Crawler) Run(seedUrl url.URL, numWorkers int, timeout time.Duration) {
	fmt.Println("Crawling", seedUrl.String())
	c.urlFrontier = make(chan *url.URL, urlFrontierBufferSize)
	c.workerResp = make(chan workerResponse, numWorkers)
	c.urlCache = newUrlCache()
	c.sitemap = newSitemap()
	c.launchWorkerPool(numWorkers, timeout)
	c.enqueueUrls([]*url.URL{&seedUrl})
	for c.tasksInProgress > 0 {
		select {
		case resp := <-c.workerResp:
			c.tasksInProgress--
			c.sitemap.addCrawledPage(resp.base, resp.linkedUrls)
			c.enqueueUrls(resp.linkedUrls)
		}
	}
	fmt.Println("Done crawling")
	c.sitemap.createDotFile()
	fmt.Println("Success")
}

func (c *Crawler) launchWorkerPool(numWorkers int, timeout time.Duration) {
	// We could launch the workers dynamically,
	// but since we expect full load on them after
	// an initial warm-up period, we can set up the pool directly
	fmt.Println("Starting up workers...")
	for i := 0; i < numWorkers; i++ {
		w := worker{
			response:    c.workerResp,
			urlFrontier: c.urlFrontier,
			timeout:     timeout,
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
