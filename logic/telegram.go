package logic

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/artemnikitin/here-tele-bot/hlp"
	"github.com/bot-api/telegram"
	"golang.org/x/net/context"
)

const (
	limitResult = 5
)

// TelegramClient describes API for bot to communicate with Telegram
type TelegramClient interface {
	SendLocationAccepted(int64)
	SendRequestForLocation(int64)
	SendUnknown(int64)
	SendWelcome(int64, string)
	SendError(int64)
	GetPlacesWithGeocoding(string, string) (*BotResult, error)
	GetPlaces(string, string) (*BotResult, error)
	SendResult(int64, *BotResult)
	SendInlineResult(string, *BotResult)
}

// TelegramMessenger
type TelegramMessenger struct {
	TelegramAPI *telegram.API
	HereAPI     *hlp.HereApiConfig
	Ctx         context.Context
	Debug       bool
}

// SendLocationAccepted answer that location was accepted
func (tm *TelegramMessenger) SendLocationAccepted(ID int64) {
	genericSend(tm, ID, locationAccepted)
}

// SendRequestForLocation sends request for location
func (tm *TelegramMessenger) SendRequestForLocation(ID int64) {
	genericSend(tm, ID, askForLocation)
}

// SendUnknown sends unknown message
func (tm *TelegramMessenger) SendUnknown(ID int64) {
	genericSend(tm, ID, unknownMessage)
}

// SendWelcome sends welcome message
func (tm *TelegramMessenger) SendWelcome(ID int64, name string) {
	genericSend(tm, ID, fmt.Sprintf(welcomeMessage, name))
}

// SendError sends error message to user
func (tm *TelegramMessenger) SendError(ID int64) {
	genericSend(tm, ID, errorHappened)
}

// GetPlacesWithGeocoding return list of places with geocoding for finding location
func (tm *TelegramMessenger) GetPlacesWithGeocoding(q, loc string) (*BotResult, error) {
	res, err := tm.HereAPI.DoGeocoding(map[string]string{
		"searchtext": loc,
		"gen":        "9",
	})
	if err != nil {
		return &BotResult{}, errors.New(err.Error())
	}
	var radius int
	switch res.Response.View[0].Result[0].MatchLevel {
	case "city":
		radius = 7000
	case "district":
		radius = 3000
	default:
		radius = 1500
	}
	lat := res.Response.View[0].Result[0].Location.DisplayPosition.Latitude
	lon := res.Response.View[0].Result[0].Location.DisplayPosition.Longitude
	return getPlacesWithRadius(tm, q, LocationToString(lat, lon), radius)
}

// GetPlaces return list of places by query when location is known
func (tm *TelegramMessenger) GetPlaces(q, loc string) (*BotResult, error) {
	return getPlacesWithRadius(tm, q, loc, 1500)
}

// SendResult send message to user with response from HERE API
func (tm *TelegramMessenger) SendResult(ID int64, results *BotResult) {
	genericSend(tm, ID, textForResponse(results))
}

// SendInlineResult send message to user with response from HERE API
func (tm *TelegramMessenger) SendInlineResult(ID string, results *BotResult) {
	inlineSend(tm, ID, results)
}

func getPlacesWithRadius(tm *TelegramMessenger, q, loc string, radius int) (*BotResult, error) {
	var wg sync.WaitGroup
	res := &BotResult{
		Location: loc,
	}
	places, err := tm.HereAPI.GetPlaces(map[string]string{
		"q":           q,
		"in":          loc + ";r=" + strconv.Itoa(radius),
		"refinements": "true",
		"tf":          "plain",
	})
	if err != nil {
		return res, errors.New(err.Error())
	}
	ch := make(chan *BotPlace, limitResult)
	count := 0
	for i := 0; i < len(places.Results.Items); i++ {
		wg.Add(1)
		go func(title, href string, dist int) {
			resp, err := tm.HereAPI.GetPlaceDetails(href)
			if err == nil {
				short, err := tm.HereAPI.ShortURL(resp.View)
				if err == nil {
					place := &BotPlace{
						Title:    title,
						Distance: dist,
						HereURL:  short.URL,
					}
					if resp.Extended.OpeningHours.IsOpen {
						place.OpeningHours = "Open now"
					} else {
						place.OpeningHours = resp.Extended.OpeningHours.Text
					}
					if len(resp.Contacts.Website) > 0 {
						place.URL = resp.Contacts.Website[0].Value
					}
					ch <- place
				}
			}
			wg.Done()
		}(places.Results.Items[i].Title, places.Results.Items[i].Href, places.Results.Items[i].Distance)
		count++
		if count == limitResult {
			break
		}
	}
	wg.Wait()
	close(ch)

	for i := range ch {
		res.Places = append(res.Places, i)
	}
	if tm.Debug {
		bytes, _ := json.Marshal(res)
		log.Println(string(bytes))
	}
	return res, nil
}

func inlineSend(tm *TelegramMessenger, ID string, results *BotResult) {
	cfg := createInlineCfg(ID)
	cfg.Results = transformResults(results)
	if tm.Debug {
		bytes, _ := json.Marshal(*cfg)
		log.Println("Inline response:", string(bytes))
	}
	if _, err := tm.TelegramAPI.AnswerInlineQuery(tm.Ctx, *cfg); err != nil {
		log.Printf("Send message error: %v", err)
	}
}

func genericSend(tm *TelegramMessenger, ID int64, text string) {
	bm := createBaseMessage(ID)
	cfg := createMessage(bm, text)
	if tm.Debug {
		bytes, _ := json.Marshal(cfg)
		log.Println("Response:", string(bytes))
	}
	if _, err := tm.TelegramAPI.SendMessage(tm.Ctx, cfg); err != nil {
		log.Printf("Send message error: %v", err)
	}
}

func createInlineCfg(ID string) *telegram.AnswerInlineQueryCfg {
	return &telegram.AnswerInlineQueryCfg{
		InlineQueryID: ID,
		CacheTime:     60,
	}
}

func transformResults(results *BotResult) []telegram.InlineQueryResult {
	var res []telegram.InlineQueryResult
	if len(results.Places) > 0 {
		count := 1
		for _, v := range results.Places {
			var buf bytes.Buffer
			buf.WriteString(v.Title)
			if v.URL != "" {
				buf.WriteString(" (")
				buf.WriteString(v.URL)
				buf.WriteString(")")
			}
			buf.WriteString("\n")
			buf.WriteString("Distance: ")
			buf.WriteString(strconv.Itoa(v.Distance))
			buf.WriteString(" m.")
			buf.WriteString("\n")
			if v.OpeningHours != "" {
				buf.WriteString(strings.Replace(v.OpeningHours, "<br/>", " ", -1))
				buf.WriteString("\n")
			}
			buf.WriteString(v.HereURL)
			res = append(res, telegram.NewInlineQueryResultArticle(strconv.Itoa(count), v.Title, buf.String()))
			count++
		}
	}
	return res
}

func createBaseMessage(ID int64) *telegram.BaseMessage {
	bc := &telegram.BaseChat{
		ID: ID,
	}
	bm := &telegram.BaseMessage{
		BaseChat: *bc,
	}
	return bm
}

func createMessage(bm *telegram.BaseMessage, text string) telegram.MessageCfg {
	cfg := &telegram.MessageCfg{
		BaseMessage:           *bm,
		Text:                  text,
		ParseMode:             telegram.MarkdownMode,
		DisableWebPagePreview: true,
	}
	return *cfg
}

func textForResponse(results *BotResult) string {
	if results.Places == nil || len(results.Places) == 0 {
		return nothingFound
	}
	var buf bytes.Buffer
	for _, v := range results.Places {
		buf.WriteString("*")
		buf.WriteString(v.Title)
		buf.WriteString("*")
		if v.URL != "" {
			buf.WriteString(" (")
			buf.WriteString(v.URL)
			buf.WriteString(")")
		}
		buf.WriteString("\n")
		buf.WriteString("Distance: ")
		buf.WriteString(strconv.Itoa(v.Distance))
		buf.WriteString(" m. ")
		buf.WriteString("\n")
		if v.OpeningHours != "" {
			buf.WriteString(strings.Replace(v.OpeningHours, "<br/>", " ", -1))
			buf.WriteString("\n")
		}
		buf.WriteString(v.HereURL)
		buf.WriteString("\n\n")
	}
	return buf.String()
}
