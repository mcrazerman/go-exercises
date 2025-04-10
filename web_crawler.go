package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type UrlCache struct {
	mu   sync.Mutex
	urls map[string]bool
}

func (uc *UrlCache) Add(key string) bool {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	if _, ok := uc.urls[key]; ok {
		return true
	} else {
		uc.urls[key] = true
		return false
	}
}

var cache UrlCache = UrlCache{urls: make(map[string]bool)}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, results chan string) {
	defer close(results)
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	if depth <= 0 {
		return
	}
	cache.Add(url)
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		results <- err.Error()
		return
	}
	results <- fmt.Sprintf("found: %s %q\n", url, body)
	temp_results := make(map[int]chan string)
	for i, u := range urls {
		if !cache.Add(u) {
			temp_results[i] = make(chan string)
			go Crawl(u, depth-1, fetcher, temp_results[i])
		}
	}

	for _, temp_res := range temp_results {
		for tr := range temp_res {
			results <- tr
		}
	}

	return
}

func main() {
	crawl_results := make(chan string)
	go Crawl("https://golang.org/", 4, fetcher, crawl_results)

	for r := range crawl_results {
		fmt.Println(r)
	}
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
