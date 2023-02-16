package crawler

import (
	"io"
	"net/http"
	"net/url"
)

type Crawler struct {
	queue   Queue[url.URL]
	visited map[url.URL]bool

	// MaxVisits limits the number of visits the crawler makes in a single run.
	// The default is zero and it means disabled.
	MaxVisits int

	// incremented at each visit
	numVisits int

	// FilterFn is an user provided function to filter which URLs to crawl
	// This function should evaluate to TRUE for URLs that SHOULD be crawled
	FilterFn func(url.URL) bool

	// VisitFn is an user provided function to execute once for each URL crawled
	VisitFn func(u url.URL, content string)
}

func NewCrawler(website string) (*Crawler, error) {
	u, err := url.Parse(website)
	if err != nil {
		return nil, err
	}

	c := Crawler{
		visited: map[url.URL]bool{},
	}

	c.Push(*u)
	return &c, nil
}

func (c *Crawler) Push(website url.URL) {
	_, ok := c.visited[website]
	if ok {
		// it's already on the list to crawl
		// don't need to put it again
		return
	}
	c.queue.Push(website)
}

func (c *Crawler) Crawl() error {
	for !c.queue.IsEmpty() && (c.MaxVisits == 0 || c.numVisits < c.MaxVisits) {
		website, err := c.queue.Pop()
		if err != nil {
			return err
		}

		err = c.Visit(website)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Crawler) Visit(u url.URL) error {
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	text, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	links, err := ExtractLinks(string(text))
	if err != nil {
		return err
	}

	for _, link := range links {
		if c.FilterFn == nil || c.FilterFn(link) {
			c.Push(link)
		}
	}

	if c.VisitFn != nil {
		c.VisitFn(u, string(text))
	}
	c.visited[u] = true

	c.numVisits++

	return nil
}
