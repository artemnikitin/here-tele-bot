package im

import (
	"github.com/patrickmn/go-cache"
)

// Config
type Config struct {
	TelegramBotKey string
	SlackToken     string
	SlackUser      string
	HereAppCode    string
	HereAppToken   string
	Cache          *cache.Cache
	Debug          bool
}
