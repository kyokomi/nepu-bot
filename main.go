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
	"github.com/kyokomi/goroku"
	"github.com/kyokomi/nepu-bot/plugins/nepubot"
	"github.com/kyokomi/slackbot"
	"github.com/kyokomi/slackbot/plugins/cron"
	"github.com/kyokomi/slackbot/plugins/kohaimage"
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
	addr, pass := goroku.GetHerokuRedisAddr()
	cronCtx := cron.NewCronContext(cron.NewRedisRepository(addr, pass, 1))
	defer cronCtx.Close()
	cronCtx.AllRefreshCron(botCtx)

	d := docomo.NewClient(apikey)
	redisRepository := NewRedisRepository()
	// add plugin
	botCtx.AddPlugin("cron", cron.NewPlugin(cronCtx))
	botCtx.AddPlugin("koha", kohaimage.NewPlugin(kohaimage.NewKohaAPI()))
	botCtx.AddPlugin("naruhodo", naruhodo.NewPlugin())
	botCtx.AddPlugin("lgtm", lgtm.NewPlugin())
	botCtx.AddPlugin("suddendeath", suddendeath.NewPlugin())
	botCtx.AddPlugin("nepu", nepubot.NewPlugin(botCtx.Plugins, d, redisRepository))

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
