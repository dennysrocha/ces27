package main

import (
	"fmt"
	"sync"
)

/*
-----------------------------------------------------------------
código consultado em: https://golang.org/src/sync/example_test.go
-----------------------------------------------------------------
*/

//var visitedUrls map[string]bool

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, visitedUrls map[string]bool) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	var mutex sync.WaitGroup //variavel de sincronismo p/ regiao critica

	if depth <= 0 {
		return
	}

	urlVisited, ok := visitedUrls[url]
	if urlVisited && ok {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	visitedUrls[url] = true //marcar pagina como visitada

	for _, u := range urls {
		mutex.Add(1) //regiao critica
		go func(u string) {
			defer mutex.Done() //terminar o mutex
			Crawl(u, depth-1, fetcher, visitedUrls)
		}(u)
	}
	mutex.Wait()
	return
}

func main() {
	visitedUrls := make(map[string]bool)
	Crawl("https://golang.org/", 4, fetcher, visitedUrls)
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
