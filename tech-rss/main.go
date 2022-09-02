package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Options struct {
	URLPath string
}

func FindAttr(
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
func FindTag[T comparable](
	node *html.Node,
	mapFunc func(n *html.Node) (t T, ok bool),
	results *[]T) {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			t, ok := mapFunc(c)
			if ok {
				*results = append(*results, t)
			}
			FindTag(c, mapFunc, results)
		}
	}
}

func ToAbsURLs(baseURL string, rssLinks []string) []string {
	results := make([]string, len(rssLinks))
	for i, rss := range rssLinks {
		u, err := url.Parse(rss)
		if err != nil {
			continue
		}
		if u.IsAbs() {
			results[i] = rss
		} else {
			results[i] = path.Join(baseURL, rss)
		}
	}
	return results
}

func rssFinder(n *html.Node) (t string, ok bool) {
	if n.DataAtom != atom.Link {
		return "", false
	}
	typeAttr, ok := FindAttr(n, "type")
	if !ok || typeAttr != "application/atom+xml" {
		return "", false
	}
	href, ok := FindAttr(n, "href")
	if !ok {
		return "", false
	}
	return href, true
}

func GetRSS(baseURL string) ([]string, error) {
	res, err := http.Get(baseURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	rssLinks := make([]string, 0)
	FindTag(doc, rssFinder, &rssLinks)
	results := ToAbsURLs(baseURL, rssLinks)
	return results, nil
}

func GetURLs(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	urls := strings.Split(string(data), "\n")
	return urls, nil
}

func run(opts *Options) error {
	urls, err := GetURLs(opts.URLPath)
	if err != nil {
		return fmt.Errorf("GetURLs failed:%w, ", err)
	}
	for _, u := range urls {
		urls, err := GetRSS(u)
		if err != nil {
			return err
		}
		fmt.Printf("%v", urls)
	}
	return nil
}

func main() {
	opts := Options{}
	flag.StringVar(&opts.URLPath, "url_path", "../data/url.txt", "text file")
	flag.Parse()
	if err := run(&opts); err != nil {
		panic(err)
	}
}
