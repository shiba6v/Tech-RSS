package main

import (
	"flag"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type Options struct {
	URLPath       string
	SlackAPIToken string
}

func run(opts *Options) error {
	urls, err := getURLs(opts.URLPath)
	api := slack.New(opts.SlackAPIToken)
	if err != nil {
		return fmt.Errorf("GetURLs failed:%w, ", err)
	}

	allChannes, err := getAllChannelNames(api)
	if err != nil {
		return nil
	}

	for _, urlString := range urls {
		u, err := url.Parse(urlString)
		if err != nil {
			return errors.WithStack(err)
		}
		channelName := createChannelName(u)
		fmt.Println(channelName)
		if !contains(allChannes, channelName) {
			rssLinks, err := getRSS(u)
			if err != nil {
				return err
			}
			// チャンネルがないなら作る。
			if err := RegisterRSSToNewChannel(api, channelName, rssLinks[0]); err != nil {
				return err
			}
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
