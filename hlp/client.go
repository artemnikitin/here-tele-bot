package hlp

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	shortenerURL = "https://maps.here.com/api/sharer/shorten?"
	searchURL    = "https://places.api.here.com/places/v1/discover/search?"
	geocoderURL  = "https://geocoder.api.here.com/6.2/search.json?"
)

type HereApiConfig struct {
	AppID    string
	AppToken string
}

func (c *HereApiConfig) ShortURL(url string) (*HereShortURLResponse, error) {
	result := &HereShortURLResponse{}
	err := execute(c, shortenerURL, map[string]string{"url": url}, result)
	return result, err
}

func (c *HereApiConfig) GetPlaces(params map[string]string) (*HerePlacesResponse, error) {
	result := &HerePlacesResponse{}
	err := execute(c, searchURL, params, result)
	return result, err
}

func (c *HereApiConfig) GetPlaceDetails(url string) (*HerePlaceDetailsResponse, error) {
	result := &HerePlaceDetailsResponse{}
	err := execute(c, url, nil, result)
	return result, err
}

func (c *HereApiConfig) DoGeocoding(params map[string]string) (*HereGeocodingResponse, error) {
	result := &HereGeocodingResponse{}
	err := execute(c, geocoderURL, params, result)
	return result, err
}

func execute(creds *HereApiConfig, url string, params map[string]string, result interface{}) error {
	req, err := http.NewRequest(http.MethodGet, createURL(creds, url, params), nil)
	if err != nil {
		return errors.New(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New(err.Error())
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New(err.Error())
	}
	err = transformBody(bytes, result)
	if err != nil {
		return errors.New(err.Error())
	}
	if resp.StatusCode != 200 {
		return errors.New("Response status code: " + strconv.Itoa(resp.StatusCode))
	}
	return nil

}

func createURL(creds *HereApiConfig, link string, params map[string]string) string {
	var res string
	if params != nil {
		var buf bytes.Buffer
		buf.WriteString(link)
		if strings.Contains(link, ".api.here.com") {
			buf.WriteString("&app_id=")
			buf.WriteString(creds.AppID)
			buf.WriteString("&app_code=")
			buf.WriteString(creds.AppToken)
		}
		for k, v := range params {
			buf.WriteString("&")
			buf.WriteString(url.QueryEscape(k))
			buf.WriteString("=")
			buf.WriteString(url.QueryEscape(v))
		}
		res = buf.String()
	} else {
		res = link
	}
	return res
}

func transformBody(bytes []byte, result interface{}) error {
	err := json.Unmarshal(bytes, result)
	return err
}
