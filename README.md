# GAWR: A crawler written in Go

GAWR is a web crawler implemented in Go using a breadth first search (BFS) algorithm.

This is roughly based on my previous crawler implementation available [here](https://github.com/danicat/spinarago). Any similarities are not a coincidence :)

## CLI Usage

First build the CLI with `go build`. Then you can crawl any website with the command `gawr`.

Please note that the default implementation uses a filter function that restricts crawling to a single domain, and the fastest crawling speed is limited to one visit per second. This will very likely change in future implementations, but the underlying crawler code can be imported as a package into your project so you are free to play with filters and limits.

```
gawr [website] [flags]

Examples:
gawr -f 1 -m 10 https://example.com

Flags:
      --config string    config file (default is $HOME/.gawr.yaml)
  -f, --frequency int    Frequency in seconds. e.g. 10 means sending one crawling request every 10 seconds. (default 10)
  -h, --help             help for gawr
  -m, --max-visits int   Maximum number of links to visit (0 = disabled)
```

## Development notes

The most important pieces of this implementation are the `VisitFn` and `FilterFn` user defined functions. `VisitFn` is called once for each URL visited and its signature is `func (u url.URL, content string)`, where u is the URL being visited and content is the full body of the page. One might write a `VisitFn` function to write this data to a key value store, for example. The default implementation on `gawr` CLI is simply a log function.

`FilterFn` on the other hand allows you to say which urls should be visited. `FilterFn` signature is `func (u url.URL) bool` and it should return true for every URL that needs to be visited. Current CLI implementation uses a simple `strings.HasPrefix` function to limit the domain being crawled, but this could easily ([?](https://xkcd.com/208/)) be replaced by a regular expression.

Additionally, in an effort to be a gentle bot, the CLI is limited to a minimum rate of one request per second, but one could simply patch the code to allow higher rates. Do it at your own risk. :)

For the same reason the code is single threaded, but maybe in the future it will be adapted to run a separate goroutine per domain.

## Contributing

If you want to contribute to this project, please raise an issue first. I'll do my best effort to triage issues and review PRs in a timely manner.

## Contact Info

For questions or comments, follow me on [Twitter](https://twitter.com/danicat83) or [Mastodon](https://hachyderm.io/danicat).