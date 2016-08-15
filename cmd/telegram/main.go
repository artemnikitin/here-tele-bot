package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/artemnikitin/here-tele-bot/im"
	c "github.com/patrickmn/go-cache"
)

var (
	debug = flag.Bool("debug", false, "Enable debug output")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	conf := &im.Config{
		TelegramBotKey: os.Getenv("BOT_TELEGRAM_KEY"),
		HereAppCode:    os.Getenv("BOT_HERE_CODE"),
		HereAppToken:   os.Getenv("BOT_HERE_TOKEN"),
		Cache:          c.New(45*time.Minute, 5*time.Minute),
		Debug:          *debug,
	}

	if conf.HereAppCode == "" || conf.HereAppToken == "" {
		fmt.Println("Set correct credentials")
		os.Exit(1)
	}

	if conf.TelegramBotKey != "" {
		im.RunTelegram(conf)
	}
}
