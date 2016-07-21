package hlp

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

var api = &HereApiConfig{
	AppID:    os.Getenv("BOT_HERE_CODE"),
	AppToken: os.Getenv("BOT_HERE_TOKEN"),
}

func TestShortURL(t *testing.T) {
	data, err := api.ShortURL("https://share.here.com/p/s-Yz1tdXNldW07aWQ9Mjc2dTMzZGItNWIwYmFkNjFjN2I5NDdmZmI4YWRmMGM1YjcwZWEwZWE7bGF0PTUyLjUyOTc1O2xvbj0xMy4zNzk4NztuPU11c2V1bStmJUMzJUJDcitOYXR1cmt1bmRlO25sYXQ9NTIuNTI5NzU7bmxvbj0xMy4zNzk4NztwaD0lMkI0OSszMC0yMDkzODU5MTtoPTE0MWI3OA")
	if err != nil {
		t.Error("Getting error: ", err.Error())
	}
	if data.Success != true {
		printStruct(data)
		t.Error("Expected to be true")
	}
}

func TestShortURLWithIncorrectURL(t *testing.T) {
	_, err := api.ShortURL("http://example.com")
	if err == nil {
		t.Error("Should return error")
	}
}

func TestGetPlaces(t *testing.T) {
	params := map[string]string{
		"q":  "asian",
		"in": "52.5303581,13.3848515;r=5000",
	}
	data, err := api.GetPlaces(params)
	if err != nil {
		t.Error("Getting error: ", err.Error())
	}
	if len(data.Results.Items) == 0 {
		printStruct(data)
		t.Error("Don't get correct response")
	}
}

func TestGeocoding(t *testing.T) {
	params := map[string]string{
		"searchtext": "mitte",
		"gen":        "9",
	}
	data, err := api.DoGeocoding(params)
	if err != nil {
		t.Error("Getting error: ", err.Error())
	}
	if len(data.Response.View) == 0 {
		printStruct(data)
		t.Error("Don't get correct response")
	}
}

func printStruct(v interface{}) {
	bytes, _ := json.Marshal(v)
	log.Println("Received:", string(bytes))
}
