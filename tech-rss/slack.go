package main

import (
	"fmt"
	"net/url"
	"strings"
	_ "unsafe"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

// /feed コマンドを打たせるのは難しい。
// 自前でRSSを巡回してメッセージをPOSTする。
// // https://qiita.com/horihiro/items/8c601c24492d87ddb742
// func SlackAPIChatCommand(api *slack.Client, channelID string, command string, args string) error {
// 	v := reflect.ValueOf(api).Elem()
// 	endpoint := *(*string)(unsafe.Pointer(v.FieldByName("endpoint").UnsafeAddr()))
// 	token := *(*string)(unsafe.Pointer(v.FieldByName("token").UnsafeAddr()))
// 	values := url.Values{
// 		"name=\"command\"": {command},
// 		"name=\"disp\"":    {command},
// 		"name=\"text\"":    {args},
// 		"name=\"channel\"": {channelID},
// 		"name=\"token\"":   {token},
// 	}
// 	fmt.Println(values)
// 	resp, err := http.PostForm(endpoint+"chat.command", values)
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}
// 	defer resp.Body.Close()
// 	respBody, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}
// 	fmt.Printf("%+v", string(respBody))
// 	return nil
// }

// チャンネルがなければ作って、generalで報告する。
func RegisterRSSToNewChannel(api *slack.Client, channelName string) error {
	_, err := api.CreateConversation(channelName, false)
	if err != nil {
		return errors.WithStack(err)
	}
	api.PostMessage("#general", slack.MsgOptionText(fmt.Sprintf("Channel #%s is created", channelName), false))
	return nil
}

func getAllChannelNames(api *slack.Client) (map[string]*slack.Channel, error) {
	params := slack.GetConversationsParameters{
		Types:           []string{"public_channel"},
		Limit:           1000,
		ExcludeArchived: true,
		Cursor:          "",
	}
	channels := make(map[string]*slack.Channel)
	for {
		resp, cursor, err := api.GetConversations(&params)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		params.Cursor = cursor
		for _, r := range resp {
			channels[r.Name] = &r
		}
		if cursor == "" {
			break
		}
	}
	return channels, nil
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
	return "r_" + channelName
}
