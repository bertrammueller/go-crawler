package crawler

import "fmt"

type Worker struct {
	// Responses from individual workers
	fetchedHtmlData chan string
	// Filtered URLs to be crawled
	urlFrontier chan string
}

func (w *Worker) Run() {
	fmt.Println("Worker running")
	for {
		select {
		case url := <-w.urlFrontier:
			fmt.Println("Fetching url", url)
		}
	}
}
