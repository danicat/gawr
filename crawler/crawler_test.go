package crawler_test

import (
	"gawr/crawler"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

type Webpage struct {
	Title   string
	Message string
	Links   []struct {
		Name string
		Link string
	}
}

func TestMain(m *testing.M) {
	templ := template.Must(template.ParseFiles("testdata/template.html"))

	sitemap := map[string]Webpage{
		"/": {
			Title:   "Home",
			Message: "Hello World",
			Links: []struct {
				Name string
				Link string
			}{
				{
					Name: "A",
					Link: "http://localhost:8080/a.html",
				},
				{
					Name: "B",
					Link: "http://localhost:8080/b.html",
				},
			},
		},
		"/a.html": {
			Title:   "Page A",
			Message: "Welcome to A",
			Links: []struct {
				Name string
				Link string
			}{
				{
					Name: "C",
					Link: "http://localhost:8080/c.html",
				},
			},
		},
		"/b.html": {
			Title:   "Page B",
			Message: "Welcome to B",
		},
		"/c.html": {
			Title:   "Page C",
			Message: "Welcome to C",
			Links: []struct {
				Name string
				Link string
			}{
				{
					Name: "A",
					Link: "http://localhost:8080/a.html",
				},
				{
					Name: "example.com",
					Link: "https://example.com/",
				},
			},
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, ok := sitemap[r.URL.RequestURI()]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := templ.Execute(w, data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	go http.ListenAndServe(":8080", nil)
	os.Exit(m.Run())
}

func TestCrawler_HasNoCycles(t *testing.T) {
	c, err := crawler.NewCrawler("http://localhost:8080")
	if err != nil {
		t.Fatal(err)
	}

	counter := map[string]int{}

	c.MaxVisits = 100
	c.VisitFn = func(u url.URL, content string) {
		counter[u.String()]++
	}

	c.FilterFn = func(u url.URL) bool {
		return strings.HasPrefix(u.String(), "http://localhost")
	}

	err = c.Crawl()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range counter {
		if v > 1 {
			t.Errorf("visit count for %v should be one, got %d", k, v)
		}
	}
}
