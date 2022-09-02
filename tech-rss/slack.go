package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/slack-go/slack"
)

func RegisterRSSToChannel(api *slack.Client, channelName string, rssLink *url.URL) error {
	// url := rssLink.String()

	// get all channels
	//
	// channel, err := api.CreateConversation(, false)
	return nil
}

func getAllChannelNames(api *slack.Client) ([]*string, error) {
	// TODO Slack API
	return nil, nil
}

func createChannelName(baseURL *url.URL) string {
	var channelName string
	if baseURL.Path == "" {
		channelName = baseURL.Host
	} else {
		channelName = fmt.Sprintf("%s__%s", baseURL.Host, baseURL.Path)
	}
	channelName = strings.Replace(channelName, "/", "_", -1)
	channelName = strings.Replace(channelName, ".", "_", -1)
	fmt.Println(channelName)
	return channelName
}
