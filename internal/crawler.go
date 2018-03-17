package crawler

import (
	"flag"
	"fmt"
)

type Crawler struct {
	// Responses from individual workers
	fetchedHtmlData chan string
	// Filtered URLs to be crawled
	urlFrontier chan string
}

func (c *Crawler) Run(domain string) {
	fmt.Println("Crawling", domain)
	c.urlFrontier = make(chan string)
	c.fetchedHtmlData = make(chan string)
	c.launchWorkerPool(1)
	c.urlFrontier <- domain
	for {
	}
}

func (c *Crawler) launchWorkerPool(numWorkers int) {
	// We could launch the workers dynamically,
	// but since we expect full load on them after
	// an initial warm-up period, we can set up the pool directly
	for i := 0; i < numWorkers; i++ {
		w := Worker{
			fetchedHtmlData: c.fetchedHtmlData,
			urlFrontier:     c.urlFrontier,
		}
		go w.Run()
	}
}

func main() {
	domain := flag.String("domain", "monzo.com", "Domain to crawl")
	flag.Parse()
	crawler := new(Crawler)
	crawler.Run(*domain)
}
