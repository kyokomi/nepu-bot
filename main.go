package main

import (
	"flag"
	"os"
	"net/http"

	"github.com/guregu/kami"
	"github.com/kyokomi/nepu-bot/bot"
	"golang.org/x/net/context"

	// init insert plugins
	_ "github.com/kyokomi/nepu-bot/plugins/nepubot"
)

func main() {
	var apikey string
	flag.StringVar(&apikey, "d", os.Getenv("DOCOMO_APIKEY"), "ドコモのAPIKEY")
	var token string
	flag.StringVar(&token, "token", os.Getenv("SLACK_BOT_TOKEN"), "SlackのBotToken")
	flag.Parse()

	ctx := context.Background()
	ctx = bot.NewDocomoClient(ctx, apikey)

	c := bot.DefaultConfig()
	c.Name = "いーすん"
	c.DocomoAPIKey = apikey
	c.SlackToken = token

	bot.WebSocketRTM(ctx, c)

	kami.Get("/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	kami.Get("/ping", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PONG"))
	})
	kami.Serve()
}
