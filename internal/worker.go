package crawl

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type worker struct {
	// Filtered URLs to be crawled
	urlFrontier chan *url.URL
	// Responses from individual workers
	fetchedURLs chan *url.URL
}

func (w *worker) Run() {
	for {
		select {
		case url := <-w.urlFrontier:
			w.doWork(url)
		}
	}
}

func (w *worker) doWork(baseUrl *url.URL) {
	fmt.Println("Fetching url", baseUrl.String())
	body, err := w.fetchWebsite(baseUrl.String())
	if err != nil {
		fmt.Println("Error trying to fetch", baseUrl.String(), ":", err.Error())
		return
	}
	fmt.Println("Received response for", baseUrl.String())
	defer body.Close()
	linkedUrls := extractUrls(body, baseUrl)
	for _, linkedUrl := range linkedUrls {
		w.fetchedURLs <- linkedUrl
	}
}

func (w *worker) fetchWebsite(url string) (io.ReadCloser, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		return nil, errors.New(fmt.Sprintf("Invalid status code %d", resp.StatusCode))
	}
	return resp.Body, err
}
