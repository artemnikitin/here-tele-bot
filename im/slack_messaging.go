package im

import (
	"github.com/artemnikitin/here-tele-bot/common"
	"github.com/artemnikitin/here-tele-bot/hlp"
	"github.com/nlopes/slack"
)

// SlackMessenger
type SlackMessenger struct {
	SlackClient *slack.Client
	HereAPI     *hlp.HereApiConfig
	Debug       bool
}

func (sm *SlackMessenger) getHLPClient() *hlp.HereApiConfig {
	return sm.HereAPI
}

func (sm *SlackMessenger) isDebug() bool {
	return sm.Debug
}

func (sm *SlackMessenger) SendError(channel string) {
	send(sm, common.ErrorHappened, channel)
}

func (sm *SlackMessenger) SendResult(result *common.BotResult, channel string) {
	send(sm, common.TextForResponse(result), channel)
}

func (sm *SlackMessenger) SendUnknown(channel string) {
	send(sm, common.SlackUnknownMessage, channel)
}

func send(sm *SlackMessenger, text, channel string) {
	params := &slack.PostMessageParameters{
		AsUser:   true,
		Markdown: true,
	}
	sm.SlackClient.PostMessage(channel, text, *params)
}
