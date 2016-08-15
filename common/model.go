package common

import "github.com/artemnikitin/here-tele-bot/hlp"

// BotResult represents response from HERE API
type BotResult struct {
	Location string
	Places   []*BotPlace
}

// BotPlace represent a single item from HERE API response
type BotPlace struct {
	Title        string
	Distance     int
	OpeningHours string
	URL          string
	HereURL      string
	IconURL      string
}

// BotInterface
type BotInterface interface {
	getHLPClient() *hlp.HereApiConfig
	isDebug() bool
}
