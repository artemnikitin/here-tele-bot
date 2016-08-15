package common

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/artemnikitin/here-tele-bot/hlp"
)

const (
	limitResult = 5
)

func GetPlacesWithGeocoding(api *hlp.HereApiConfig, q, loc string) (*BotResult, error) {
	res, err := api.DoGeocoding(map[string]string{
		"searchtext": loc,
		"gen":        "9",
	})
	if err != nil {
		return &BotResult{}, errors.New(err.Error())
	}
	var radius int
	var lat, lon float64
	if len(res.Response.View) == 0 {
		return &BotResult{}, errors.New("Empty geocoding response")
	}
	switch res.Response.View[0].Result[0].MatchLevel {
	case "city":
		radius = 7000
	case "district":
		radius = 3000
	default:
		radius = 1500
	}
	lat = res.Response.View[0].Result[0].Location.DisplayPosition.Latitude
	lon = res.Response.View[0].Result[0].Location.DisplayPosition.Longitude
	return GetPlacesWithRadius(api, q, LocationToString(lat, lon), radius)
}

func GetPlacesWithRadius(api *hlp.HereApiConfig, q, loc string, radius int) (*BotResult, error) {
	var wg sync.WaitGroup
	res := &BotResult{
		Location: loc,
	}
	places, err := api.GetPlaces(map[string]string{
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
			resp, err := api.GetPlaceDetails(href)
			if err == nil {
				short, err := api.ShortURL(resp.View)
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
					if resp.Icon != "" {
						place.IconURL = resp.Icon
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
	/*if api.isDebug() {
		bytes, _ := json.Marshal(res)
		log.Println(string(bytes))
	}*/
	return res, nil
}

// TextForResponse transforms result from HERE API to message
func TextForResponse(results *BotResult) string {
	if results.Places == nil || len(results.Places) == 0 {
		return NothingFound
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
