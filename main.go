package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/guregu/kami"
	"github.com/kyokomi/slackbot"
	"github.com/kyokomi/slackbot/plugins"
	"golang.org/x/net/context"

	// init insert bot.plugins
	//	_ "github.com/kyokomi/slackbot/plugins/echo"

	"github.com/kyokomi/nepu-bot/plugins/nepubot"
	"github.com/kyokomi/slackbot/plugins/cron"
	_ "github.com/kyokomi/slackbot/plugins/lgtm"
	_ "github.com/kyokomi/slackbot/plugins/suddendeath"
	_ "github.com/kyokomi/slackbot/plugins/tiqav"
)

//go:generate ego -package main

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

	cron.Setup()
	defer cron.Stop()

	slackbot.WebSocketRTM(ctx, c)

	kami.Get("/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		IndexTmpl(w, plugins.GetPlugins())
	})
	kami.Get("/ping", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PONG"))
	})
	kami.Serve()
}
