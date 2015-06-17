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

	// 6 docomo APIとかで対話
	"github.com/kyokomi/nepu-bot/plugins/nepubot"
	// 5 LGTM画像ランダム
	_ "github.com/kyokomi/slackbot/plugins/lgtm"
	// 4 突然死のやつ
	_ "github.com/kyokomi/slackbot/plugins/suddendeath"
	// 3 tiqavで画像検索
	_ "github.com/kyokomi/slackbot/plugins/tiqav"
	// 2 twitterで画像検索
	_ "github.com/kyokomi/nepu-bot/plugins/twitter"
	// 1 cronの設定
	"github.com/kyokomi/slackbot/plugins/cron"
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
