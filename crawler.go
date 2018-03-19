package main

import (
	"flag"
	"fmt"
	"net/url"
	"time"

	crawl "github.com/bertgit/go-crawler/internal/crawl"
)

func main() {
	flagUrl := flag.String("url", "http://www.monzo.com", "url to crawl")
	flagNumWorkers := flag.Int("workers", 100, "Number of workers")
	flagTimeout := flag.Int("timeout", 10, "Timeout in seconds per http request")
	flag.Parse()
	crawler := new(crawl.Crawler)
	baseUrl, err := url.Parse(*flagUrl)
	if err != nil || !baseUrl.IsAbs() {
		fmt.Println("Invalid URL", *flagUrl)
		return
	}
	crawler.Run(*baseUrl, *flagNumWorkers, time.Duration(*flagTimeout)*time.Second)
}
