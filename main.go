package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/guregu/kami"
	"github.com/kyokomi/slackbot"
	"golang.org/x/net/context"
	"github.com/kyokomi/slackbot/plugins"

	// init insert bot.plugins
//	_ "github.com/kyokomi/slackbot/plugins/echo"
	"github.com/kyokomi/nepu-bot/plugins/nepubot"
)

func main() {
	var apikey string
	flag.StringVar(&apikey, "d", os.Getenv("DOCOMO_APIKEY"), "ドコモのAPIKEY")
	var token string
	flag.StringVar(&token, "token", os.Getenv("SLACK_BOT_TOKEN"), "SlackのBotToken")
	flag.Parse()

	ctx := plugins.Context()
	ctx = nepubot.NewContext(ctx, apikey)

	c := slackbot.DefaultConfig()
	c.Name = "いーすん"
	c.SlackToken = token

	slackbot.WebSocketRTM(ctx, c)

	kami.Get("/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	kami.Get("/ping", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PONG"))
	})
	kami.Serve()
}
