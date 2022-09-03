package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func findAttr(
	node *html.Node,
	key string) (attr string, ok bool) {
	for _, v := range node.Attr {
		if v.Key == key {
			return v.Val, true
		}
	}
	return "", false
}

// https://golang.hateblo.jp/entry/golang-net-html
func findTag[T comparable](
	node *html.Node,
	mapFunc func(n *html.Node) (t T, ok bool),
	results *[]T) {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			t, ok := mapFunc(c)
			if ok {
				*results = append(*results, t)
			}
			findTag(c, mapFunc, results)
		}
	}
}

func toAbsURLs(baseURL *url.URL, rssURLs []*url.URL) []*url.URL {
	results := make([]*url.URL, len(rssURLs))
	for i, rss := range rssURLs {
		if rss.IsAbs() {
			results[i] = rss
		} else {
			rss.Scheme = baseURL.Scheme
			rss.Host = baseURL.Host
			results[i] = rss
		}
	}
	return results
}

func rssFinder(n *html.Node) (t *url.URL, ok bool) {
	if n.DataAtom != atom.Link {
		return nil, false
	}
	typeAttr, ok := findAttr(n, "type")
	isFeed := typeAttr == "application/atom+xml" ||
		typeAttr == "application/rss+xml"
	if !ok || !isFeed {
		return nil, false
	}
	href, ok := findAttr(n, "href")
	if !ok {
		return nil, false
	}
	hrefURL, err := url.Parse(href)
	if err != nil {
		return nil, false
	}
	return hrefURL, true
}

func getRSS(baseURL *url.URL) ([]*url.URL, error) {
	res, err := http.Get(baseURL.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()
	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	rssLinks := make([]*url.URL, 0)
	findTag(doc, rssFinder, &rssLinks)
	results := toAbsURLs(baseURL, rssLinks)
	if len(rssLinks) == 0 {
		return nil, fmt.Errorf("the page does not have RSS: %v", baseURL)
	}
	return results, nil
}
