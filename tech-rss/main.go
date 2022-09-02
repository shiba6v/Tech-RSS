package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Options struct {
	URLPath       string
	SlackAPIToken string
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

func ToAbsURLs(baseURL *url.URL, rssURLs []*url.URL) []*url.URL {
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
	typeAttr, ok := FindAttr(n, "type")
	isFeed := typeAttr == "application/atom+xml" ||
		typeAttr == "application/rss+xml"
	if !ok || !isFeed {
		return nil, false
	}
	href, ok := FindAttr(n, "href")
	if !ok {
		return nil, false
	}
	hrefURL, err := url.Parse(href)
	if err != nil {
		return nil, false
	}
	return hrefURL, true
}

func GetRSS(baseURL *url.URL) ([]*url.URL, error) {
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
	FindTag(doc, rssFinder, &rssLinks)
	results := ToAbsURLs(baseURL, rssLinks)
	return results, nil
}

func GetURLs(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	urls := strings.Split(strings.TrimSpace(string(data)), "\n")
	for i, u := range urls {
		urls[i] = strings.TrimRight(u, "/")
	}
	return urls, nil
}

func RegisterRSSToChannel(api *slack.Client, baseURL *url.URL, rssLink *url.URL) error {
	// url := rssLink.String()
	var channelName string
	if baseURL.Path == "" {
		channelName = baseURL.Host
	} else {
		channelName = fmt.Sprintf("%s__%s", baseURL.Host, baseURL.Path)
	}
	channelName = strings.Replace(channelName, "/", "_", -1)
	channelName = strings.Replace(channelName, ".", "_", -1)
	fmt.Println(channelName)
	// channel, err := api.CreateConversation(, false)
	return nil
}

func run(opts *Options) error {
	urls, err := GetURLs(opts.URLPath)
	api := slack.New(opts.SlackAPIToken)
	if err != nil {
		return fmt.Errorf("GetURLs failed:%w, ", err)
	}
	for _, urlString := range urls {
		u, err := url.Parse(urlString)
		if err != nil {
			return errors.WithStack(err)
		}
		rssLinks, err := GetRSS(u)
		if err != nil {
			return err
		}
		if len(rssLinks) == 0 {
			return fmt.Errorf("the page does not have RSS: %s", u)
		}
		if err := RegisterRSSToChannel(api, u, rssLinks[0]); err != nil {
			return err
		}
		// fmt.Printf("%v\n", rssLinks)
	}
	return nil
}

func main() {
	opts := Options{}
	flag.StringVar(&opts.URLPath, "url_path", "../data/url.txt", "text file")
	flag.StringVar(&opts.SlackAPIToken, "slack_api_token", "", "slack api token")
	flag.Parse()
	if err := run(&opts); err != nil {
		fmt.Printf("Error:%+v", err)
	}
}
