package im

import (
	"fmt"
	"log"
	"strings"

	"github.com/artemnikitin/here-tele-bot/common"
	"github.com/artemnikitin/here-tele-bot/hlp"
	"github.com/nlopes/slack"
)

//TODO: Use cache ?
// RunSlack runs bot
func RunSlack(config *Config) {
	slackID := ""
	api := slack.New(config.SlackToken)
	api.SetDebug(config.Debug)
	sm := &SlackMessenger{
		SlackClient: api,
		HereAPI: &hlp.HereApiConfig{
			AppID:    config.HereAppCode,
			AppToken: config.HereAppToken,
		},
		Debug: config.Debug,
	}

	users, err := sm.SlackClient.GetUsers()
	if err != nil {
		log.Println("Can't get list of Slack users:", err)
		return
	}
	for _, user := range users {
		if user.Name == config.SlackUser {
			slackID = user.ID
			break
		}
	}

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Println("Slack event received")
			data, ok := msg.Data.(*slack.MessageEvent)
			if ok {
				ID := data.User
				text := data.Text
				if data.SubMessage != nil {
					ID = data.SubMessage.User
					text = data.SubMessage.Text
				}
				if ID != slackID {
					if strings.Contains(text, "<!here|@here>.:") || strings.Contains(text, "<!here>.:") {
						if strings.Contains(text, " in ") {
							text = common.ClearSlackMessage(text)
							q, loc := common.SplitQueryAndLocation(text)
							places, err := common.GetPlacesWithGeocoding(sm.HereAPI, q, loc)
							if err != nil {
								sm.SendError(data.Channel)
								return
							}
							sm.SendResult(places, data.Channel)
						} else {
							sm.SendUnknown(data.Channel)
						}
					}
				}
			}

		}
	}
}
