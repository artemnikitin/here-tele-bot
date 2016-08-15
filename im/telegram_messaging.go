package im

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/artemnikitin/here-tele-bot/common"
	"github.com/artemnikitin/here-tele-bot/hlp"
	"github.com/bot-api/telegram"
	"golang.org/x/net/context"
)

// TelegramMessenger
type TelegramMessenger struct {
	TelegramAPI *telegram.API
	HereAPI     *hlp.HereApiConfig
	Ctx         context.Context
	Debug       bool
}

func (tm *TelegramMessenger) getHLPClient() *hlp.HereApiConfig {
	return tm.HereAPI
}

func (tm *TelegramMessenger) isDebug() bool {
	return tm.Debug
}

// SendLocationAccepted answer that location was accepted
func (tm *TelegramMessenger) SendLocationAccepted(ID int64) {
	genericSend(tm, ID, common.LocationAccepted)
}

// SendRequestForLocation sends request for location
func (tm *TelegramMessenger) SendRequestForLocation(ID int64) {
	genericSend(tm, ID, common.TelegramAskForLocation)
}

// SendUnknown sends unknown message
func (tm *TelegramMessenger) SendUnknown(ID int64) {
	genericSend(tm, ID, common.TelegramUnknownMessage)
}

// SendWelcome sends welcome message
func (tm *TelegramMessenger) SendWelcome(ID int64, name string) {
	genericSend(tm, ID, fmt.Sprintf(common.TelegramWelcomeMessage, name))
}

// SendError sends error message to user
func (tm *TelegramMessenger) SendError(ID int64) {
	genericSend(tm, ID, common.ErrorHappened)
}

// GetPlaces return list of places by query when location is known
func (tm *TelegramMessenger) GetPlaces(q, loc string) (*common.BotResult, error) {
	return common.GetPlacesWithRadius(tm.HereAPI, q, loc, 1500)
}

// SendResult send message to user with response from HERE API
func (tm *TelegramMessenger) SendResult(ID int64, results *common.BotResult) {
	genericSend(tm, ID, common.TextForResponse(results))
}

// SendInlineResult send message to user with response from HERE API
func (tm *TelegramMessenger) SendInlineResult(ID string, results *common.BotResult) {
	inlineSend(tm, ID, results)
}

func inlineSend(tm *TelegramMessenger, ID string, results *common.BotResult) {
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

func transformResults(results *common.BotResult) []telegram.InlineQueryResult {
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
			item := telegram.NewInlineQueryResultArticle(strconv.Itoa(count), v.Title, buf.String())
			if v.URL != "" {
				item.URL = v.URL
			}
			buf.Reset()
			buf.WriteString("Distance: ")
			buf.WriteString(strconv.Itoa(v.Distance))
			buf.WriteString(" m.")
			item.Description = buf.String()
			item.HideURL = false
			if v.IconURL != "" {
				item.ThumbURL = v.IconURL
			}
			res = append(res, item)
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
