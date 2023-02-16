package main

import (
	"gawr/crawler"
	"log"
	"net/url"
	"strings"
)

func main() {
	startURL := "https://example.com"

	c, err := crawler.NewCrawler(startURL)
	if err != nil {
		log.Fatalf("fatal error crawling %v: %v", startURL, err)
	}

	c.MaxVisits = 3
	c.VisitFn = func(u url.URL, content string) {
		log.Printf("visiting %#v", u.String())
	}

	c.FilterFn = func(u url.URL) bool {
		return strings.HasPrefix(u.String(), startURL)
	}

	err = c.Crawl()
	if err != nil {
		log.Printf("error crawling: %v", err)
	}
}
