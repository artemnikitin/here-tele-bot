package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/artemnikitin/here-tele-bot/hlp"
	"github.com/artemnikitin/here-tele-bot/logic"
	"github.com/bot-api/telegram"
	c "github.com/patrickmn/go-cache"
	"golang.org/x/net/context"
)

var (
	debug = flag.Bool("debug", false, "Enable debug output")

	botKey       = os.Getenv("BOT_KEY")
	hereAppCode  = os.Getenv("BOT_HERE_CODE")
	hereAppToken = os.Getenv("BOT_HERE_TOKEN")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	if botKey == "" || hereAppToken == "" || hereAppCode == "" {
		fmt.Println("Set correct credentials")
		os.Exit(1)
	}
	api := telegram.New(botKey)
	api.Debug(*debug)
	cache := c.New(45*time.Minute, 5*time.Minute)
	ctx, cancel := context.WithCancel(context.Background())
	tm := &logic.TelegramMessenger{
		TelegramAPI: api,
		HereAPI: &hlp.HereApiConfig{
			AppID:    hereAppCode,
			AppToken: hereAppToken,
		},
		Ctx:   ctx,
		Debug: *debug,
	}
	defer cancel()

	if user, err := api.GetMe(ctx); err != nil {
		log.Panic(err)
	} else {
		log.Println("Bot started:", user)
	}

	updatesCh := make(chan telegram.Update)

	go telegram.GetUpdates(ctx, api, telegram.UpdateCfg{
		Timeout: 5, // Timeout in seconds for long polling.
		Offset:  0, // Start with the oldest update
	}, updatesCh)

	for update := range updatesCh {
		if update.InlineQuery != nil {
			processInlineQuery(tm, update)
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
				processNormalQuery(cache, tm, update)
			}
		} else {
			tm.SendUnknown(update.Chat().ID)
		}

	}

}

func processInlineQuery(tm *logic.TelegramMessenger, update telegram.Update) {
	query := update.InlineQuery
	if strings.Contains(query.Query, " in ") {
		q, loc := logic.SplitQueryAndLocation(query.Query)
		places, err := tm.GetPlacesWithGeocoding(q, loc)
		if err != nil {
			tm.SendError(update.From().ID)
			return
		}
		tm.SendInlineResult(update.InlineQuery.ID, places)
	} else {
		if query.Location == nil {
			return
		}
		loc := logic.LocationToString(query.Location.Latitude, query.Location.Longitude)
		places, err := tm.GetPlaces(query.Query, loc)
		if err != nil {
			tm.SendError(update.From().ID)
			return
		}
		tm.SendInlineResult(update.InlineQuery.ID, places)
	}
}

func processNormalQuery(cache *c.Cache, tm *logic.TelegramMessenger, update telegram.Update) {
	if update.Message.Location != nil {
		lat := update.Message.Location.Latitude
		lon := update.Message.Location.Longitude
		loc := logic.LocationToString(lat, lon)
		cache.Set(update.From().Username, loc, c.DefaultExpiration)
		tm.SendLocationAccepted(update.Chat().ID)
		return
	}
	if strings.Contains(update.Message.Text, " in ") {
		q, loc := logic.SplitQueryAndLocation(update.Message.Text)
		places, err := tm.GetPlacesWithGeocoding(q, loc)
		if err != nil {
			tm.SendError(update.Chat().ID)
			return
		}
		cache.Set(update.From().Username, places.Location, c.DefaultExpiration)
		tm.SendResult(update.Chat().ID, places)
	} else {
		coord, ok := cache.Get(update.From().Username)
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
