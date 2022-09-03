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
	api := slack.New(opts.SlackAPIToken)
	db, err := setupDB()
	if err != nil {
		return err
	}

	urls, err := getURLs(opts.URLPath)

	if err != nil {
		return fmt.Errorf("GetURLs failed:%w, ", err)
	}

	allChannes, err := getAllChannelNames(api)
	if err != nil {
		return err
	}
	fmt.Println(allChannes)

	for _, urlString := range urls {
		u, err := url.Parse(urlString)
		if err != nil {
			return errors.WithStack(err)
		}
		channelName := createChannelName(u)
		fmt.Println(channelName)
		// チャンネルがないなら作る。
		if _, ok := allChannes[channelName]; !ok {
			if err := RegisterRSSToNewChannel(api, channelName); err != nil {
				return err
			}
		}

		rssLinks, err := getRSS(u)
		if err != nil {
			return err
		}
		CreateSubscription(db, rssLinks[0].String(), channelName)

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
