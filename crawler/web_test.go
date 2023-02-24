package crawler_test

import (
	"gawr/crawler"
	"net/url"
	"testing"
)

func TestExtractLinks(t *testing.T) {
	input := `
<html>
<head>
</head>
<body>
<p><a href="https://www.example.com">example</a></p>
<br>
<p><a href="#something">section</a></p>
<br>
<p><a href="/relative.html">relative</a></p>
`
	src, _ := url.Parse("https://foo.com")

	expected := []string{
		"https://www.example.com",
		"#something",
		"/relative.html",
	}

	results, err := crawler.ExtractLinks(src, input)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != len(expected) {
		t.Errorf("expected %d results, got %d", len(expected), len(results))
	}

	for i, link := range results {
		if link.String() != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], link)
		}
	}
}
