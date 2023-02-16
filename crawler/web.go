package crawler

import (
	"log"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// ExtractLinks takes the HTML text of a webpage and returns all links in HREF tags.
func ExtractLinks(text string) ([]url.URL, error) {
	reader := strings.NewReader(text)
	tokenizer := html.NewTokenizer(reader)

	var links []url.URL

	for tt := tokenizer.Next(); tt != html.ErrorToken; tt = tokenizer.Next() {
		if tt == html.StartTagToken {
			t := tokenizer.Token()

			if t.Data == "a" {
				for _, attr := range t.Attr {
					if attr.Key == "href" {
						u, err := url.Parse(attr.Val)
						if err != nil {
							log.Printf("error parsing href tag with value %v: %v", attr.Val, err)
							continue
						}

						links = append(links, *u)
					}
				}
			}
		}
	}

	return links, nil
}
