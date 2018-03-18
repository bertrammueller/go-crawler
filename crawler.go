package main

import (
	"flag"
	"fmt"
	"net/url"

	crawl "github.com/bertgit/crawler/internal/crawl"
)

func main() {
	flagUrl := flag.String("url", "http://www.monzo.com", "url to crawl")
	flagNumWorkers := flag.Int("workers", 100, "Number of workers")
	flag.Parse()
	crawler := new(crawl.Crawler)
	baseUrl, err := url.Parse(*flagUrl)
	if err != nil || !baseUrl.IsAbs() {
		fmt.Println("Invalid URL", *flagUrl)
		return
	}
	crawler.Run(*baseUrl, *flagNumWorkers)
}
