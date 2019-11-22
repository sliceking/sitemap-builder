package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/svwielga4/link-parser"
)

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url you want to build a sitemap for")
	depthFlag := flag.Int("depth", 3, "the depth you want the bfs to search for links")
	flag.Parse()
	pages := bfs(*urlFlag, *depthFlag)
	// pages := get(*urlFlag)
	for _, href := range pages {
		fmt.Println(href)
	}
}

func bfs(urlStr string, depth int) []string {
	// create a map string:empty struct to track what we've seen
	seen := make(map[string]struct{})
	// make a similar map for our queue of pages
	var q map[string]struct{}
	// and another one for our 'next queue' starting with the base
	nq := map[string]struct{}{
		urlStr: struct{}{},
	}
	// loop through however deep you want the bfs to go
	for i := 0; i <= depth; i++ {
		// make the current queue the next queue, and the next queue a new queue
		q, nq = nq, make(map[string]struct{})
		for page := range q {
			// if its already seen, skip
			if _, ok := seen[page]; ok {
				continue
			}
			// mark it as seen
			seen[page] = struct{}{}
			// get all the links off of that page and put them into the next queue
			for _, link := range get(page) {
				nq[link] = struct{}{}
			}
		}
	}
	var ret []string
	// change the seen map into a slice
	for k := range seen {
		ret = append(ret, k)
	}
	return ret
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	base := baseURL.String()
	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(body io.Reader, base string) []string {
	links, _ := link.Parse(body)
	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}
