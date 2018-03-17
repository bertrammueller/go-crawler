package crawl

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type worker struct {
	// Filtered URLs to be crawled
	urlFrontier chan url.URL
	// Responses from individual workers
	fetchedURLs chan url.URL
}

func (w *worker) Run() {
	fmt.Println("Worker running")
	for {
		select {
		case url := <-w.urlFrontier:
			w.doWork(url)
		}
	}
}

func (w *worker) doWork(curUrl url.URL) {
	fmt.Println("Fetching url", curUrl)
	html, err := w.fetchHtml(curUrl.String())
	if err != nil {
		fmt.Println("Error trying to fetch", curUrl, ":", err.Error())
		return
	}
	fmt.Println("Received html for", curUrl, "size", len(html))
	w.fetchedURLs <- curUrl
}

func (w *worker) fetchHtml(url string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New(fmt.Sprintf("Invalid status code %d", resp.StatusCode))
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	return string(bytes), err
}
