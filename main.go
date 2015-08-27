package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/braintree/manners"
	"github.com/labstack/echo"
	"github.com/zenazn/goji/bind"

	"github.com/kyokomi/go-docomo/docomo"
	"github.com/kyokomi/nepu-bot/plugins/nepubot"
	"github.com/kyokomi/slackbot"
	"github.com/kyokomi/slackbot/plugins/cron"
	"github.com/kyokomi/slackbot/plugins/lgtm"
	"github.com/kyokomi/slackbot/plugins/naruhodo"
	"github.com/kyokomi/slackbot/plugins/suddendeath"
)

//go:generate ego -package main

func init() {
	bind.WithFlag()
	if fl := log.Flags(); fl&log.Ltime != 0 {
		log.SetFlags(fl | log.Lmicroseconds)
	}
}

func main() {
	listener := bind.Default()
	log.Println("Starting on", listener.Addr())

	var apikey string
	flag.StringVar(&apikey, "d", os.Getenv("DOCOMO_APIKEY"), "ドコモのAPIKEY")
	var token string
	flag.StringVar(&token, "token", os.Getenv("SLACK_BOT_TOKEN"), "SlackのBotToken")
	if !flag.Parsed() {
		flag.Parse()
	}

	botCtx, err := slackbot.NewBotContext(token)
	if err != nil {
		panic(err)
	}
	// cronを設定
	cronCtx := cron.NewCronContext(cron.NewHerokuRedisRepository())
	defer cronCtx.Close()
	cronCtx.AllRefreshCron(botCtx)

	d := docomo.NewClient(apikey)

	// add plugin
	nepuBotPlugin := nepubot.Plugin{Docomo: d, Plugins: botCtx.Plugins}
	botCtx.AddPlugin("nepu", &nepuBotPlugin)
	botCtx.AddPlugin("cron", cron.Plugin{CronContext: cronCtx})
	botCtx.AddPlugin("naruhodo", naruhodo.Plugin{})
	botCtx.AddPlugin("lgtm", lgtm.Plugin{})
	botCtx.AddPlugin("suddendeath", suddendeath.Plugin{})

	// start
	botCtx.WebSocketRTM()

	// herokuで動くように

	e := echo.New()
	e.Get("/", func(w http.ResponseWriter, r *http.Request) {
		IndexTmpl(w, botCtx.Plugins.GetPlugins())
	})
	e.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PONG"))
	})

	manners.Serve(listener, e.Router())
}
