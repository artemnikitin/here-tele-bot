package im

import (
	"strings"

	"github.com/artemnikitin/here-tele-bot/common"
	"github.com/artemnikitin/here-tele-bot/hlp"
	"github.com/bot-api/telegram"
	"github.com/patrickmn/go-cache"
	"golang.org/x/net/context"
)

// RunTelegram run bot for Telegram
func RunTelegram(config *Config) {
	api := telegram.New(config.TelegramBotKey)
	api.Debug(config.Debug)
	ctx, cancel := context.WithCancel(context.Background())
	tm := &TelegramMessenger{
		TelegramAPI: api,
		HereAPI: &hlp.HereApiConfig{
			AppID:    config.HereAppCode,
			AppToken: config.HereAppToken,
		},
		Ctx:   ctx,
		Debug: config.Debug,
	}
	defer cancel()

	updatesCh := make(chan telegram.Update)

	go telegram.GetUpdates(ctx, api, telegram.UpdateCfg{
		Timeout: 5, // Timeout in seconds for long polling.
		Offset:  0, // Start with the oldest update
	}, updatesCh)

	for update := range updatesCh {
		if update.InlineQuery != nil {
			processInlineQuery(tm, update)
			continue
		}
		if update.HasMessage() {
			if isChoiceFromInlineResult(update) {
				continue
			}
			if isStartCommand(update) {
				var user string
				if update.From().FirstName != "" {
					user = update.From().FirstName
				} else {
					user = update.From().Username
				}
				tm.SendWelcome(update.Chat().ID, user)
				continue
			}
			if update.Message.Text != "" || update.Message.Location != nil {
				processNormalQuery(config.Cache, tm, update)
				continue
			}
		} else {
			tm.SendUnknown(update.Chat().ID)
		}

	}
}

func processInlineQuery(tm *TelegramMessenger, update telegram.Update) {
	query := update.InlineQuery
	if w, ok := common.IsQueryCorrect(query.Query); ok {
		q, loc := common.SplitQueryAndLocation(query.Query, w)
		places, err := common.GetPlacesWithGeocoding(tm.HereAPI, q, loc)
		if err != nil {
			tm.SendError(update.From().ID)
			return
		}
		tm.SendInlineResult(update.InlineQuery.ID, places)
	} else {
		if query.Location == nil {
			return
		}
		loc := common.LocationToString(query.Location.Latitude, query.Location.Longitude)
		places, err := tm.GetPlaces(query.Query, loc)
		if err != nil {
			tm.SendError(update.From().ID)
			return
		}
		tm.SendInlineResult(update.InlineQuery.ID, places)
	}
}

func processNormalQuery(c *cache.Cache, tm *TelegramMessenger, update telegram.Update) {
	if update.Message.Location != nil {
		lat := update.Message.Location.Latitude
		lon := update.Message.Location.Longitude
		loc := common.LocationToString(lat, lon)
		c.Set(update.From().Username, loc, cache.DefaultExpiration)
		tm.SendLocationAccepted(update.Chat().ID)
		return
	}
	if w, ok := common.IsQueryCorrect(update.Message.Text); ok {
		q, loc := common.SplitQueryAndLocation(update.Message.Text, w)
		places, err := common.GetPlacesWithGeocoding(tm.HereAPI, q, loc)
		if err != nil {
			tm.SendError(update.Chat().ID)
			return
		}
		c.Set(update.From().Username, places.Location, cache.DefaultExpiration)
		tm.SendResult(update.Chat().ID, places)
	} else {
		coord, ok := c.Get(update.From().Username)
		if ok {
			v, ok := coord.(string)
			if ok {
				places, err := tm.GetPlaces(update.Message.Text, v)
				if err != nil {
					tm.SendError(update.Chat().ID)
					return
				}
				tm.SendResult(update.Chat().ID, places)
			} else {
				tm.SendError(update.Chat().ID)
			}
		} else {
			tm.SendRequestForLocation(update.Chat().ID)
		}
	}
}

func isChoiceFromInlineResult(msg telegram.Update) bool {
	result := false
	if msg.HasMessage() {
		if len(msg.Message.Entities) > 0 && strings.Contains(msg.Message.Text, "Distance: ") {
			result = true
		}
	}
	return result
}

func isStartCommand(msg telegram.Update) bool {
	result := false
	if msg.HasMessage() && msg.Message.Entities != nil {
		if msg.Message.Entities[0].Type == "bot_command" && msg.Message.Text == "/start" {
			result = true
		}
	}
	return result
}
